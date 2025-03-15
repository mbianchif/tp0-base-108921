import signal
import logging
from server.common.protocol import BetSockListener, BetSockStream
from server.common.utils import Bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        self._listener = BetSockListener.bind("", port, listen_backlog)
        self._shutdown = False

        def term_handler(_signum, _stacktrace):
            self._shutdown = True

        signal.signal(signal.SIGTERM, term_handler)

    def run(self):
        while not self._shutdown:
            stream = self._accept_new_connection()
            try:
                self._handle_client_connection(stream)
            finally:
                stream.close()

        self._listener.close()

    def _handle_client_connection(self, client_sock: BetSockStream):
        try:
            msg = client_sock.recv()
            bet = Bet(
                msg._agency,
                msg._name,
                msg._surname,
                msg._id,
                msg._birthdate,
                msg._number,
            )

            store_bets([bet])
            logging.info(f"action: apuesta_almacenada | result: success | dni: ${msg._id} | numero: ${msg._number}")
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")

    def _accept_new_connection(self) -> BetSockStream:
        logging.info("action: accept_connections | result: in_progress")
        conn, addr = self._listener.accept()
        logging.info(f"action: accept_connections | result: success | ip: {addr[0]}")
        return conn
