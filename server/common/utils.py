import csv
import datetime
from typing import Iterator

DELIMITER = ","

""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574


class Bet:
    """A lottery bet registry."""

    def __init__(
        self,
        agency: str,
        first_name: str,
        last_name: str,
        document: str,
        birthdate: str,
        number: str,
    ):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

    @classmethod
    def from_bytes(cls, data: bytearray):
        return cls(*data.decode().split(DELIMITER))


def has_won(bet: Bet) -> bool:
    """Checks whether a bet won the prize or not."""
    return bet.number == LOTTERY_WINNER_NUMBER


def store_bets(bets: list[Bet]) -> None:
    """
    Persist the information of each bet in the STORAGE_FILEPATH file.
    Not thread-safe/process-safe.
    """
    with open(STORAGE_FILEPATH, "a+") as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow(
                [
                    bet.agency,
                    bet.first_name,
                    bet.last_name,
                    bet.document,
                    bet.birthdate,
                    bet.number,
                ]
            )


def load_bets() -> Iterator[Bet]:
    """
    Loads the information all the bets in the STORAGE_FILEPATH file.
    Not thread-safe/process-safe.
    """
    with open(STORAGE_FILEPATH, "r") as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])
