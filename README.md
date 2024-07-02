## Local setup

```bash
ddev start
ddev composer install
ddev import-db --file=dump.sql.gz
```

Run console command via:
```bash
ddev craft my-module/product/create
```

User: admin, password: password123

## How was this created

```bash
ddev config \
    --project-type=craftcms \
    --docroot=web \
    --create-docroot \
    --php-version="8.2" \
    --database="mysql:8.0" \
    --nodejs-version="20" && \
  ddev start -y && \
  ddev composer create -y --no-scripts craftcms/craft && \
  ddev craft install/craft \
    --username=admin \
    --password=password123 \
    --email=admin@example.com \
    --site-name=Testsite \
    --language=en \
    --site-url='$DDEV_PRIMARY_URL' && \
  echo 'Nice, ready to launch!' && \
  ddev launch

# Active craft pro trial in plugin store

# add generator:
ddev composer require craftcms/generator --dev

# install commerce:
ddev composer require "craftcms/commerce:^5.0.11.1" -w && ddev exec php craft plugin/install commerce
# --> activate pro version in dashboard

# add timber
ddev composer require "verbb/timber:^2.0.1" -w && ddev exec php craft plugin/install timber
```

```bash
craft make module
# Selected
# Module ID: (kebab-case) my-module
# Base module class name: (PascalCase) [Module] MyModule
# Module location: modules
# Should the module be loaded during app initialization? (yes|no) [no]:yes
```

- Add product type general
- Add console command (modules)

Thanks to https://github.com/craftcms/commerce/discussions/3458