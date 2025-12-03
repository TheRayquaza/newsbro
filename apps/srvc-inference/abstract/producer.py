import logging
from typing import List

import pydantic
from kafka import KafkaProducer

from abstract.message import InferenceCommand


class InferenceProducerConfig(pydantic.BaseModel):
    kafka_bootstrap_servers: str
    kafka_producer_topic: str


class InferenceProducer:
    def __init__(self, logger: logging.Logger, config: InferenceProducerConfig) -> None:
        self.config = config
        self.producer = KafkaProducer(
            bootstrap_servers=self.config.kafka_bootstrap_servers,
            value_serializer=lambda v: v.json().encode("utf-8"),
            acks="all",
            linger_ms=5,
            retries=3,
            max_in_flight_requests_per_connection=1,
        )
        self.logger = logger

    def produce(self, data: List[InferenceCommand]) -> None:
        for item in data:
            try:
                self.producer.send(self.config.kafka_producer_topic, value=item)
            except Exception as e:
                self.logger.error(f"Error producing message: {e}", exc_info=True)
        self.producer.flush()
