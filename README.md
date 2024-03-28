# citadel-corp/segokuning-social-app

I'm Commander Shepard, and This Is My Favorite 'Social App Service' on the Citadel

🔍 Tested with
[this k6 script](https://github.com/nandanugg/SocialMediaTestCases).

📝 Documentation - TBA.

🎵 Songs to test by - [playlist](https://open.spotify.com/album/1oVSp3g7ULNAHzFtdBvHEd?si=IVw3cdo6RUKDOdb1gYCJKQ).

## Getting Started

These instructions will give you a copy of the project up and running on
your local machine for development and testing purposes. See deployment
for notes on deploying the project on a live system.

### Prerequisites

Requirements for the software and other tools to run and test
- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [k6](https://k6.io/docs/get-started/installation/) - to test

Note that we use [AWS S3 service](https://aws.amazon.com/s3/) to upload image,
[setup your own](https://docs.aws.amazon.com/AmazonS3/latest/userguide/GetStartedWithS3.html) S3 bucket if you want to test uploading image.

### Migrate the database

After [setting up your database locally](https://www.postgresql.org/docs/current/tutorial-createdb.html),
run this to migrate our database structure
```
$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2
```

### Running the service

A step by step series of that tell you how to get a development
environment running

Create file named .env and fill it with, example:
```
DB_HOST = localhost
DB_PORT = 5432
DB_USERNAME = root
DB_PASSWORD = pass12345
DB_NAME = segokuning
DB_PARAMS = "&sslmode=disable"
BCRYPT_SALT = 12
JWT_SECRET = ${JWT_SECRET}
S3_ID = ${S3_ID}
S3_SECRET_KEY = ${S3_SECRET_KEY}
S3_BUCKET_NAME = ${S3_BUCKET_NAME}
S3_REGION = ${S3_REGION}
ENV = local
```

Run the service

    go run cmd/main.go

Now you can run the service. Service is running on port `8080`.

## Endpoints
- User
    - Register - `POST /v1/user/register`
    - Login - `POST /v1/user/login`
    - Link Email - `POST /v1/user/link`
    - Link Email - `POST /v1/user/link/email`
    - Link Phone - `POST /v1/user/link/phone`
    - Update - `PATCH /v1/user`
- Friends
    - Add - `POST /v1/friend`
    - Remove - `DELETE /v1/friend`
    - List - `GET /v1/friend`
- Post
    - Create - `POST /v1/post`
    - List - `GET /v1/post`
    - Comment - `POST /v1/post/comment`
- Image
    - Upload - `POST /v1/image`

## Running the tests

First, run our service.

After [installing k6](https://k6.io/docs/get-started/installation/), 
[clone this test script](https://github.com/nandanugg/SocialMediaTestCases) outside of this project directory, and then run:
```
$ BASE_URL=http://localhost:8080 k6 run script.js
```

To run with certain number of iterations and virtual users:
```
$ BASE_URL=http://localhost:8080 k6 run -i NUMBER_OF_ITERATIONS --vus NUMBER_OF_VIRTUAL_USERS script.js
```

To run with real world cases:
```
BASE_URL=http://localhost:8080 REAL_WORLD_CASE=true k6 run -i NUMBER_OF_ITERATIONS --vus NUMBER_OF_VIRTUAL_USERS script.js
```

## Authors

The [Citadel Corp](https://github.com/citadel-corp) team:
  - [**TheSagab**](https://github.com/TheSagab)
  - [**Faye**](https://github.com/farolinar)

## License

This project is licensed under the [MIT License](https://github.com/citadel-corp/segokuning-social-app?tab=MIT-1-ov-file) - see the [LICENSE](https://github.com/citadel-corp/segokuning-social-app/blob/main/LICENSE) file for
details

## Acknowledgments

  - The Ramadhan ProjectSprint organizer and members
