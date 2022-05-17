from locust import HttpUser, task

import string
import json
import random
import decimal

# Generate random items
class ItemLoad(HttpUser):
    def generateItemName(self):
        return ''.join(random.choices(string.ascii_lowercase + string.digits, k=32))

    def generateItemValue(self):
        return str(decimal.Decimal(random.randrange(100, 10000))/100)

    @task
    def createItem(self):
        headers = {"Content-Type": "application/json"}
        data = {"item_name": self.generateItemName(), "item_value": self.generateItemValue()}

        self.client.post("/items", data=json.dumps(data), headers=headers)

