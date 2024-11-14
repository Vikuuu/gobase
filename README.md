# gobase
automatic schema detection that creates the migrations for your code

![gobase-img](https://github.com/Vikuuu/gobase/assets/img/gobase-img.png) 

## Goal
This is the goal for the gobase v0.0.1 

Input code:
```go
package main

import (
    "time"
)

type users struct {
    ID        int 
    CreatedAt time.Time 
    UpdatedAt time.Time 
    Name      string
    Email     string
    IsMember  bool
}
```

Output code:
```sql
CREATE TABLE users (
    id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT,
    email TEXT,
    is_member BOOLEAN
);
```
