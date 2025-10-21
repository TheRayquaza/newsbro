import json
import logging
import os
from threading import Thread
from typing import Any, Dict

import pydantic
from kafka import KafkaConsumer


class InferenceConsumerConfig(pydantic.BaseModel):
    kafka_bootstrap_servers: str = os.getenv(
        "KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"
    )
    kafka_consumer_topic: str = os.getenv("KAFKA_CONSUMER_TOPIC", "")
    auto_offset_reset: str = "earliest"
    kafka_consumer_group: str = os.getenv(
        "KAFKA_CONSUMER_GROUP", "inference-consumer-group"
    )
    auto_commit: bool = True
    batch_size: int = int(os.getenv("BATCH_SIZE", "50"))
    batch_interval: int = int(os.getenv("BATCH_INTERVAL", "2"))


class InferenceConsumer:
    def __init__(
        self,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        config: dict,
        health_hook=None,
        process_hook=None,
        *args,
        **kwargs,
    ):
        self.config = config
        self.logger = logger
        self.consumer = KafkaConsumer(
            consumer_config.kafka_consumer_topic,
            bootstrap_servers=consumer_config.kafka_bootstrap_servers,
            auto_offset_reset=consumer_config.auto_offset_reset,
            enable_auto_commit=consumer_config.auto_commit,
            group_id=consumer_config.kafka_consumer_group,
            value_deserializer=lambda x: json.loads(x.decode("utf-8")),
        )
        self.state: Dict[str, Any] = {}
        self.consumer_config = consumer_config
        self.health_hook = health_hook if health_hook else (lambda: True)
        self.process_hook = process_hook if process_hook else (lambda batch: None)

    def health(self) -> bool:
        if not self.health_hook():
            return False
        topics = self.consumer.topics()
        if not topics:
            return False
        return True

    def run_impl(self) -> None:
        batch = []
        try:
            for message in self.consumer:
                article = message.value
                batch.append(article)
                if len(batch) >= self.consumer_config.batch_size:
                    Thread(target=self.process_hook, args=(batch,), daemon=True).start()
                    batch = []
        except Exception as e:
            self.logger.error(f"Error in consumer loop: {e}", exc_info=True)

    def run(self) -> None:
        self.logger.info("Starting periodic batch processor thread.")
        return Thread(target=self.run_impl, daemon=True).start()
