from locust import HttpUser, task, events
from locust.exception import RescheduleTask

import string
import json
import random
import requests

# Generate games
# A game consists of 100 players. Only 1 winner randomly selected from those players
#
# Matchmaking is random list of players that are not playing
#
# To achieve this
# A locust user 'GameMatch' will start off by creating a "game"
# Then, pre-selecting a subset of users, and set a current_game attribute for those players.
# Once done, after a period of time, a winner is randomly selected.


# TODO: Matchmaking should ideally be handled by Agones. Once done, Locust test would convert to testing Agones match-making
# Create and close game matches
class GameMatch(HttpUser):

    @task
    def createGame(self):
        headers = {"Content-Type": "application/json"}

        # Create the game
        # TODO: Make number of players configurable
        # data = {"numPlayers": 10}
        res = self.client.post("/games/create", headers=headers)

        # TODO: Store the response into memory to be used to close the game later, to avoid a call to the DB

    @task
    def closeGame(self):
        # Get a game that's currently open, then close it
        headers = {"Content-Type": "application/json"}
        with self.client.get("/games/open", headers=headers, catch_response=True) as response:
            try:
                data = {"gameUUID": response.json()["gameUUID"]}
                self.client.put("/games/close", data=json.dumps(data), headers=headers)
            except json.JSONDecodeError:
                response.failure("Response could not be decoded as JSON")
            except KeyError:
                response.failure("Response did not contain expected key 'playerUUID'")


