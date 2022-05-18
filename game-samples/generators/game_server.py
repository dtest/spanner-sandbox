from locust import HttpUser, task

import json
import random

# Generate player load with 5:1 reads to write
class GameLoad(HttpUser):
    def generateAmount(self):
        return str(round(random.uniform(1.01, 49.99), 2))

    @task
    def acquireMoney(self):
        headers = {"Content-Type": "application/json"}

        # Get a random player that's part of a game, and update balance
        with self.client.get("/players", headers=headers, catch_response=True) as response:
            try:
                data = {"playerUUID": response.json()["playerUUID"], "amount": self.generateAmount(), "source": "loot"}
                self.client.put("/players/balance", data=json.dumps(data), headers=headers)
            except json.JSONDecodeError:
                response.failure("Response could not be decoded as JSON")
            except KeyError:
                response.failure("Response did not contain expected key 'playerUUID'")
