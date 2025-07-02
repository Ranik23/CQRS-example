package domain



var (	// OutboxMessageKey is the key used to store outbox messages in the cache.
	OrderCreatedKey = "order_created"
	OrderUpdatedKey = "order_updated"
	OrderDeletedKey = "order_deleted"
)