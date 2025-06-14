CREATE TABLE IF NOT EXISTS goods (
    Id UInt64,
    ProjectId UInt64,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    EventTime DateTime DEFAULT now(),
    INDEX idx_ProjectId ProjectId TYPE minmax GRANULARITY 1,
    INDEX idx_Name Name TYPE minmax GRANULARITY 1
) ENGINE = MergeTree
ORDER BY Id;