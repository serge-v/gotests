#include <sys/socket.h>
#include <sys/types.h>
#include <sys/epoll.h>
#include <netdb.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <memory.h>
#include <fcntl.h>
#include <errno.h>
#include <pthread.h>

#define MAXEVENTS 1000

static int
create_and_bind(char *port)
{
	struct addrinfo hints;
	struct addrinfo *result, *rp;
	int s, sfd;

	memset(&hints, 0, sizeof(struct addrinfo));
	hints.ai_family = AF_UNSPEC;     /* Return IPv4 and IPv6 choices */
	hints.ai_socktype = SOCK_STREAM; /* We want a TCP socket */
	hints.ai_flags = AI_PASSIVE;     /* All interfaces */

	s = getaddrinfo(NULL, port, &hints, &result);
	if (s != 0)
	{
		fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(s));
		return -1;
	}

	for (rp = result; rp != NULL; rp = rp->ai_next)
	{
		sfd = socket(rp->ai_family, rp->ai_socktype, rp->ai_protocol);
		if (sfd == -1)
			continue;

		s = bind(sfd, rp->ai_addr, rp->ai_addrlen);
		if (s == 0)
		{
			/* We managed to bind successfully! */
			break;
		}

		close(sfd);
	}

	if (rp == NULL)
	{
		fprintf(stderr, "Could not bind\n");
		return -1;
	}

	freeaddrinfo(result);

	int one = 1;

	if (setsockopt(sfd, SOL_SOCKET, SO_REUSEADDR, (char *)&one, sizeof(one)) == -1)
	{
		fprintf(stderr, "Could not setsockopt\n");
		return -1;
	}
	
	return sfd;
}

static int
make_socket_non_blocking(int sfd)
{
	int flags, s;

	flags = fcntl(sfd, F_GETFL, 0);
	if (flags == -1)
	{
		perror("fcntl");
		return -1;
	}

	flags |= O_NONBLOCK;
	s = fcntl(sfd, F_SETFL, flags);
	if (s == -1)
	{
		perror("fcntl");
		return -1;
	}

	return 0;
}

const char resp[] = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 1\r\n\r\n\n";

struct params
{
	int num;
	int sfd;
};

void
accept_connection(int sfd, int efd, int num)
{
	int s;
	struct epoll_event event;

	/* We have a notification on the listening socket, which
	   means one or more incoming connections. */
	while (1)
	{
		struct sockaddr in_addr;
		socklen_t in_len;
		int infd;
		char hbuf[NI_MAXHOST], sbuf[NI_MAXSERV];

		in_len = sizeof in_addr;
		infd = accept(sfd, &in_addr, &in_len);
		if (infd == -1)
		{
			if ((errno == EAGAIN) ||
				(errno == EWOULDBLOCK))
			{
				/* We have processed all incoming
				   connections. */
				break;
			}
			else
			{
				perror("accept");
				break;
			}
		}

		s = getnameinfo(&in_addr, in_len,
				hbuf, sizeof hbuf,
				sbuf, sizeof sbuf,
				NI_NUMERICHOST | NI_NUMERICSERV);
		if (s == 0)
		{
			printf("%d: accept %d "
			       "(host=%s, port=%s)\n", num, infd, hbuf, sbuf);
		}

		/* Make the incoming socket non-blocking and add it to the
		   list of fds to monitor. */
		s = make_socket_non_blocking(infd);
		if (s == -1)
			abort();

		event.data.fd = infd;
		event.events = EPOLLIN | EPOLLET;
		s = epoll_ctl(efd, EPOLL_CTL_ADD, infd, &event);
		if (s == -1)
		{
			perror("epoll_ctl");
			abort();
		}
	}
}

void*
threadproc(void* d)
{
	int efd;
	struct epoll_event event;
	struct epoll_event *events;
	
	struct params* pars = (struct params*)d;
	int sfd = pars->sfd;

	efd = epoll_create(1000);
	if (efd == -1)
	{
		perror("epoll_create");
		abort();
	}

	event.data.fd = sfd;
	event.events = EPOLLIN | EPOLLET;
	int s = epoll_ctl(efd, EPOLL_CTL_ADD, sfd, &event);
	if (s == -1)
	{
		perror("epoll_ctl");
		abort();
	}

	/* Buffer where events are returned */
	events = calloc(MAXEVENTS, sizeof event);

	/* The event loop */
	while (1)
	{
		int n, i, readers = 0;

		n = epoll_wait(efd, events, MAXEVENTS, -1);
		for (i = 0; i < n; i++)
		{
			if ((events[i].events & EPOLLERR) ||
			        (events[i].events & EPOLLHUP) ||
			        (!(events[i].events & EPOLLIN)))
			{
				fprintf(stderr, "epoll error\n");
				close(events[i].data.fd);
				events[i].data.fd = 0;
			}
			else if (sfd != events[i].data.fd)
				readers++;
		}

		while (readers > 0)
		{
			for (i = 0; i < n; i++)
			{
				if (events[i].data.fd == 0 || sfd == events[i].data.fd)
					continue;
			
				int done = 0;

				ssize_t count;
				char buf[4096];

				count = read(events[i].data.fd, buf, sizeof buf);
				if (count == -1)
				{
					/// end of data
					if (errno != EAGAIN)
					{
						perror("read");
						done = 1;
					}
//					readers--;
//					events[i].data.fd = 0;
					readers = 0; // go to epoll_wait
					break;
				}
				else if (count == 0)
				{
					// end of file
					printf("%d: closed conn %d\n", pars->num, events[i].data.fd);
					/* Closing the descriptor will make epoll remove it
					   from the set of descriptors which are monitored. */
					close(events[i].data.fd);
					readers--;
					events[i].data.fd = 0;
					continue;
				}

				s = write(events[i].data.fd, resp, sizeof(resp)-1);

				if (s == -1)
				{
					perror("write resp");
					abort();
				}
			}
		}

		// accept
		for (i = 0; i < n; i++)
		{
			if (sfd == events[i].data.fd)
				accept_connection(sfd, efd, pars->num);
		}

	}

	free(events);

	close(sfd);

	return EXIT_SUCCESS;
}

const int THREADS = 8;

int
main(int argc, char *argv[])
{
	int sfd, s, i;
	int efd;
	struct epoll_event event;
	struct epoll_event *events;

	if (argc != 2)
	{
		fprintf(stderr, "Usage: %s [port]\n", argv[0]);
		exit(EXIT_FAILURE);
	}

	sfd = create_and_bind(argv[1]);
	if (sfd == -1)
		abort();

	s = make_socket_non_blocking(sfd);
	if (s == -1)
		abort();

	s = listen(sfd, SOMAXCONN);
	if (s == -1)
	{
		perror("listen");
		abort();
	}
	
	struct params params[THREADS];
	pthread_t threads[THREADS];
	
	for (i = 0; i < THREADS; i++)
	{
		params[i].num = i;
		params[i].sfd = sfd;
		int rc = pthread_create(&threads[i], NULL, threadproc, &params[i]);
		if (rc < 0)
		{
			perror("pthread_create");
			abort();
		}
	}

	for (i = 0; i < THREADS; i++)
	{
		int rc = pthread_join(threads[i], NULL);
		if (rc < 0)
		{
			perror("pthread_join");
			abort();
		}
	}
}
