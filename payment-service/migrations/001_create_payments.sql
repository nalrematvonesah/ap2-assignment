CREATE TABLE payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL
);