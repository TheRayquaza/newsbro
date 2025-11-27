import abc
import json
import logging
import time
from concurrent.futures import ThreadPoolExecutor
from threading import Event, Lock, Thread
from typing import Any, List

import pydantic
from kafka import KafkaConsumer


class InferenceConsumerConfig(pydantic.BaseModel):
    kafka_bootstrap_servers: str
    kafka_consumer_topic: str
    auto_offset_reset: str = "earliest"
    kafka_consumer_group: str
    auto_commit: bool = True
    batch_size: int
    batch_interval: int
    session_timeout_ms: int = 60000
    heartbeat_interval_ms: int = 10000
    max_poll_interval_ms: int = 300000
    thread_pool_size: int


class InferenceConsumer(abc.ABC):
    def __init__(
        self,
        logger: logging.Logger,
        consumer_config: InferenceConsumerConfig,
        dict_to_msg=None,
        *args,
        **kwargs,
    ):
        self.logger = logger
        self.thread: None | Thread = None
        self.dict_to_msg = dict_to_msg
        self.consumer = KafkaConsumer(
            consumer_config.kafka_consumer_topic,
            bootstrap_servers=consumer_config.kafka_bootstrap_servers,
            auto_offset_reset=consumer_config.auto_offset_reset,
            enable_auto_commit=consumer_config.auto_commit,
            group_id=consumer_config.kafka_consumer_group,
            value_deserializer=lambda x: json.loads(x.decode("utf-8")),
            auto_commit_interval_ms=5000,
            session_timeout_ms=consumer_config.session_timeout_ms,
            heartbeat_interval_ms=consumer_config.heartbeat_interval_ms,
            max_poll_interval_ms=consumer_config.max_poll_interval_ms,
        )
        self.consumer_config = consumer_config
        self.thread_pool = ThreadPoolExecutor(
            max_workers=consumer_config.thread_pool_size
        )
        self.batch: List[Any] = []
        self.batch_lock = Lock()
        self.last_process_time = time.time()
        self.shutdown_event = Event()

    def health(self) -> bool:
        try:

            if self.consumer is None:
                return False

            if self.consumer._closed:
                return False

            if self.thread is None or not self.thread.is_alive():
                return False

            return True

        except Exception:
            return False

    @abc.abstractmethod
    def process(self, batch: List[Any]):
        pass

    def _should_process_batch(self) -> bool:
        """Check if batch should be processed based on size or time. MUST BE CALLED with batch_lock held."""
        if len(self.batch) >= self.consumer_config.batch_size:
            return True

        if len(self.batch) > 0:
            time_elapsed = time.time() - self.last_process_time
            if time_elapsed >= self.consumer_config.batch_interval:
                return True

        return False

    def _process_batch(self) -> None:
        """Process current batch and reset. MUST BE CALLED with batch_lock held."""
        if not self.batch:
            return
        batch_to_process = self.batch.copy()
        self.batch = []
        self.last_process_time = time.time()

        try:
            self.logger.info(f"Processing batch of {len(batch_to_process)} messages")
            self.thread_pool.submit(self.process, batch_to_process)
        except Exception as e:
            self.logger.error(f"Error launching process thread: {e}", exc_info=True)

    def run_impl(self) -> None:
        """Main consumer loop with smart batching."""
        self.logger.info(
            f"Starting consumer with batch_size={self.consumer_config.batch_size}, "
            f"batch_interval={self.consumer_config.batch_interval}s"
        )

        try:
            while not self.shutdown_event.is_set():
                try:
                    for message in self.consumer:
                        article = message.value
                        if self.dict_to_msg:
                            article = self.dict_to_msg(article)

                        with self.batch_lock:
                            self.batch.append(article)
                            if len(self.batch) >= self.consumer_config.batch_size:
                                self._process_batch()

                        # Check time-based processing outside the message loop
                        with self.batch_lock:
                            if self._should_process_batch():
                                self._process_batch()

                except StopIteration:
                    with self.batch_lock:
                        if self._should_process_batch():
                            self._process_batch()
                except Exception as e:
                    self.logger.error(f"Error in consumer loop: {e}", exc_info=True)
        finally:
            with self.batch_lock:
                if self.batch:
                    self.logger.info("Processing remaining batch on shutdown")
                    self._process_batch()

    def run(self) -> None:
        """Start the consumer thread."""
        self.logger.info("Starting consumer thread")
        self.thread_pool.submit(self.run_impl)

    def shutdown(self) -> None:
        """Gracefully shutdown the consumer."""
        self.logger.info("Shutting down consumer")
        self.shutdown_event.set()
        self.consumer.close()
        self.thread_pool.shutdown(wait=True)
