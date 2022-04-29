# Backend diorama.id
## Requirements
- Go 1.17
- Postgresql 14
- `.env` file
## How to Run
- You can access the deployed API on http://34.101.123.15:8080/
### Local
- Install all the requirements
- Create a postgresql database named using `diorama.sql` in the db directory
- Run `go build && ./diorama` in the src directory (Note: This is a Bash based command, change it to fit your shell)
### Docker
- Run `docker-compose up` on the root directory

## Documentation
- [API Documentation](https://documenter.getpostman.com/view/19661864/UVkjwJ76)