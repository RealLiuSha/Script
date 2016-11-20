import socket
# Getting a random free tcp port in python using sockets

def get_free_tcp_port():
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(('', 0))
    addr, port = tcp.getsockname()
    tcp.close()
    return port


def get_free_tcp_address():
    tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp.bind(('', 0))
    host, port = tcp.getsockname()
    tcp.close()
    return 'tcp://{host}:{port}'.format(**locals())


def get_free_ports(howMany=10):
    """Return a list of n free port numbers on localhost"""
    results = []
    sockets = []
    for x in range(howMany):
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.bind(('localhost', 0))
        # work out what the actual port number it's bound to is
        addr, port = s.getsockname()
        results.append(port)
        sockets.append(s)

    for s in sockets:
        s.close()

    return results

print(get_free_tcp_port())
print(get_free_tcp_address())
print(get_free_ports())
