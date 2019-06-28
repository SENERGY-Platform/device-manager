- REST Create/Update/Delete -> CQRS Event
- REST Get -> aggregate REST to CQRS-Query-Services  

# API

- `POST|GET /device-types`
    - GET Parameter: 
        - limit
        - offset
        - sort (name.asc)		
    - example: `/device-types?limit=10&offset=0&sort=name.asc`
- `PUT|GET|DELETE  /device-types/:id`

- `POST|GET /devices`
    - GET Parameter: 
        - limit
        - offset
        - sort (name.asc)		
    - example: `/devices?limit=10&offset=0&sort=name.asc`
- `PUT|GET|DELETE  	/devices/:id`


-----------
**create and delete only for admins:**
- `POST|GET /protocols`
    - GET Parameter: 
        - limit
        - offset
        - sort (name.asc)		
    - example: `/protocols?limit=10&offset=0&sort=name.asc`
- `GET|DELETE  		/protocols/:id`

- `POST|GET /serializations`
    - GET Parameter: 
        - limit
        - offset
        - sort (name.asc)		
    - example: `/serializations?limit=10&offset=0&sort=name.asc`
- `GET|DELETE  		/serializations/:id`


