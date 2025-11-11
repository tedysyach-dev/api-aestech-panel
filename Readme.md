## How to use this template
### clone this project
```bash
git clone git@github.com:tedysyach-dev/boilerplate-backend-go.git
```
### delete template history
```bash
rm -rf .git
```
### init git 
```bash
git init
git add .
git commit -m "Init from boilerplate"
git remote add origin https://github.com/username/api-customer.git
git push -u origin main
```

---
## How to add migration
### create migration
```bash
make migrate-new n="migaration_name"
```
### Migration Up
```bash
make migarte-up
```
### Migration Down
```bash
make migarte-down
```