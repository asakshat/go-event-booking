# Event Booking

This is the starting template for the collaborative project "Event Booking".

## Authors

- **Denis Monseur** - _Front-End_ - [Github Link](https://github.com/DenisMonseur)
- **Dora Nacinovic** - _Front-End_ - [Github Link](https://github.com/DoraNac)
- **Sakshat Adhikari** - _Back-End_ - [Github Link](https://github.com/asakshat)

## Project Overview

The "Event Booking" project is designed to provide a platform for users to book events.

## APIs

```
URL :- https://eventbooking-go-9c6c8d14446d.herokuapp.com/
```

### Auth Requests

> Post /api/signup
>
> ```json
> {
> 	"username": "retrosax",
> 	"email": "retrosax@gmail.com",
> 	"password": "MegaMan123!"
> }
> ```

> Post /api/login
>
> ```json
> {
> 	"email": "retrosax@gmail.com",
> 	"password": "MegaMan123!"
> }
> ```

> Post /api/logout
>
> ```
> Logs out user
> ```

> Post api/logged
>
> ```
> Get the user details of the currently logged in user
> ```

### Event Routes

> Post /api/event/create # (with form-data)
>
> ```
> {
> title:Gopher
> description:Gopher description
> venue:BeCode
> date:20/09/2024
> time:18:00
> location:rude de canterlaan , 16
> price:29.99
> }
> ```

> PUT /api/event/update/:event_id #(with form data)
>
> ```
> {
> title:Gopherrrrss
> description:Gopher descriptionsss
> venue:BeCodesss
> date:20/09/2023
> time:19:00
> location:rude de canterlaan , 19
> price:19.99
> }
> ```

> Delete /api/event/delete/:event_id
>
> ```
> Deletes event ( not permanent)
> ```

> PATCH /api/event/undo-delete/:event_id
>
> ```
> Undo the deleted event
> ```

> Delete /api/event/delete_perm/:event_id
>
> ```
> Deletes event permanently
> ```

> GET /api/event/:event_id
>
> ```
> Get event by event id
> ```

> GET /api/event
>
> ```
> Get all events
> ```
