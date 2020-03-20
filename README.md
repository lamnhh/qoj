# QHH Online Judge

An online judge that uses problem structure similar to that of Themis

Demo video [here](https://www.youtube.com/watch?v=o2DUFkfa8ms)

# Requirement

- A PostgreSQL server running
- `go >= 1.11` (any version that support Go Modules)
- `node`
- `yarn`

# Usage

## Create a .env file

View `.env-sample` to see what `.env` needs to contain.

## Build Javascript

```
$ yarn install.
$ yarn build.
```

## Start server

```
$ go run main.go  // client app, for regular users
$ go run main-admin.go  // admin app, for managing problems and contests
```

The client app will be hosted on `localhost:3000` and the admin app on `localhost:3001`.

For now, you have to register on the client app and manually modify admin privilege to be able to access the admin app.

# To-Do

- [x] Admin dashboard (for managing problems, contests, test data).
- [x] Add responsive CSS.
- [ ] Improve judger (`./timeout`).
- [ ] Docker for deployment.

# License

MIT License

Copyright (c) 2020 Lam Nguyen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
