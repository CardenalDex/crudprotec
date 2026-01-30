## 🚀 CRUDProtec: Technical test for Protec/Promass job

Welcome to **CRUDProtec**, a robust Go-based microservice designed with **Clean Architecture** principles. This system handles businesses, merchants, and transactions with a precision-first approach (using cents for money) and an automated audit trail that captures every move.


## Known bugs
> this is just a random technic test, so if by some reason you find it, and by some reason is usefull to you, know that I `wont` maintain or fix it 
* **When creating a bussines remember data type in of commisssion is FLOAT, i wasnt able to fix the swagger file, and by default it just puts a "0" but it must be `0.0`**

---

## 🛠 Features

* **Clean Architecture:** Strict separation between Entities, Use Cases, Adapters, and Infrastructure.
* **CGO-Powered SQLite:** high-performance internal storage using GORM.
* **Automated Audit:** Every mutation (create/update/delete) is logged via UseCase side-effects.
* **Money Safety:** All calculations are handled in **Integer Cents/Basis Points** to avoid floating-point nightmares.
* **Live Documentation:** Fully interactive Swagger UI generated automatically on build.

---

## 🏗 How to Run

Running the entire stack (including DB initialization and Swagger generation) is distilled into a single sequence:

```bash
docker-compose build && docker-compose up

```
> That **should** be the case, but by some reason i didnt get to run it in a single cmd so....

First run 
```bash
docker-compose build

```
then 
```bash
docker-compose up

```
I dont know why but my system dindt recognice the `&&`

**What happens during the boot process?**

1. **Swagger Generation:** The `swaggger` tool scans the code and updates the documentation.
2. **CGO Compilation:** The Go binary is compiled with C-extensions for SQLite.
3. **DB Migration:** GORM detects the internal `.db` file and automatically builds your tables.
4. **Ready for Traffic:** The server starts listening on port `8080`.

**You can change settings using the ENV variables in the `docker-compose.yaml` file** 

---

## 📖 API Documentation

Once the container is spinning, you can explore, test, and break the API through the interactive Swagger interface:

🔗 **[http://localhost:8080/swagger/index.html](https://www.google.com/search?q=http://localhost:8080/swagger/index.html)**

> **Note:** Use this to register your first **Business** (getting a UUID), then a **Merchant** (getting a UUID), and finally process a **Transaction**.

## System handles ID as **UUID** not as integers
Why? in your *Definiciones de estructura sugeridas* logs use an UID, so I thniked to standarice all the project.
**I definitilly know** that thats the worst idea possible, but since this isnt a PROD ready project, i wanted to try it  

---

## 🕵️ The Audit Trail (Logs)

In this system, logs aren't just for debugging; they are a first-class citizen of the business logic.

### How it works

Unlike standard application logs that just print to the console, our **Audit Logs** are triggered inside the **UseCase (Service) layer**.

* **When:** Every time a `Create`, `Update`, or `Delete` method is called.
* **Where:** Stored in the `logs` table inside the SQLite database.
* **What:** We capture the `Action`, the `ResourceID` (UUID of the affected entity), and a timestamp.
* **Why:** This provides a non-repudiable history of system changes. If a business commission is updated, you'll see exactly when and which ID was affected.

**Retrieve logs via:** `GET /api/v1/audit` or `GET /api/v1/audit/{resource_id}`.
But is **EASIER** to do all directly from swagger

---

## 📂 Project Structure

* `cmd/app/`: Entry point (Main.go).
* `internal/entitys/`: Pure business models (Domain).
* `internal/usecase/`: Business rules and Service layer (The "Brain").
* `internal/adapter/`: Implementation details (GORM, Gin Handlers).
* `internal/adapter/repository/`: The SQLite persistence layer.

---

### 🧪 Quick Test

After running the project, try creating a business with a 5.5% commission:
**I didnt tested this curl, pls use swagger, if needed i can do all the testing and provide you with a postman file**

```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/admin/businesses' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{ "commission_percentage": 5.5 }'

```

---

This is my first time writing a readme, and i always wanted to make it pretty like those real projects