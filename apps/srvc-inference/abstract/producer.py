import logging
import os
from typing import List

import pydantic
from kafka import KafkaProducer

from abstract.message import InferenceCommand


class InferenceProducerConfig(pydantic.BaseModel):
    kafka_bootstrap_servers: str = os.getenv(
        "KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"
    )
    kafka_producer_topic: str = os.getenv("KAFKA_PRODUCER_TOPIC", "")


class InferenceProducer:
    def __init__(self, logger: logging.Logger, config: InferenceProducerConfig) -> None:
        self.config = config
        self.producer = KafkaProducer(
            bootstrap_servers=self.config.kafka_bootstrap_servers,
            value_serializer=lambda v: v.json().encode("utf-8"),
        )
        self.logger = logger

    def produce(self, data: List[InferenceCommand]) -> None:
        for item in data:
            try:
                self.producer.send(self.config.kafka_producer_topic, value=item)
            except Exception as e:
                self.logger.error(f"Error producing message: {e}", exc_info=True)
        self.producer.flush()
