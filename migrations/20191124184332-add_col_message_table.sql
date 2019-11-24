
-- +migrate Up
ALTER TABLE `messages` ADD `available_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
    -- NOT NULL would set enum to take the first value => 'created'
ALTER TABLE `messages` ADD `status` enum ('created', 'in_transit', 'queued', 'requeued', 'processed') NOT NULL;

-- +migrate Down

ALTER TABLE `messages` DROP `available_at`;
ALTER TABLE `messages` DROP `status`;
