CREATE KEYSPACE IF NOT EXISTS chat_app
WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 3};

USE chat_app;

CREATE TABLE IF NOT EXISTS messages (
    chat_id UUID,
    message_id UUID,
    sender_id UUID,
    content TEXT,
    timestamp TIMESTAMP,
    PRIMARY KEY (chat_id, timestamp, message_id)
) WITH CLUSTERING ORDER BY (timestamp DESC);
