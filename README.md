## Local setup

```bash
ddev start && ddev composer install && ddev import-db --file=dump.sql.gz
```

User: admin, password: password123

## Problems

### A) Custom required field not validated (WIP)

Run console command via:
```bash
ddev craft my-module/product/create
```

Not yet activated, see source code: 

- [ ] Custom required field for products (product type) is not validated by `Craft::$app->elements->saveElement($product)`

### B) GraphQL deadlock issues for concurrent mutations (WIP)

Run graphql test via a go(lang) script on your local laptop, three workers in parallel:

```bash
cd test-graphl-mutation/
# run script
./graphql_mutation 

# build it after changes
go build -o graphql_mutation main.go     
```

results in

```
GraphQL Errors: [map[debugMessage:SQLSTATE[40001]: Serialization failure: 1213 Deadlock found when trying to get lock; try restarting transaction extensions:map[category:internal] message:Internal server error trace:[map[call:PDOStatement::execute() file:/var/www/html/vendor/yiisoft/yii2/db/Command.php line:1320] map[call:yii\db\Command::internalExecute('INSERT INTO `entries_authors` (`entryId`, `authorId`, `sortOrder`) VALUES (268, 1, 1)') file:/var/www/html/vendor/yiisoft/yii2/db/Command.php line:1120] map[call:yii\db\Command::execute() file:/var/www/html/vendor/craftcms/cms/src/helpers/Db.php line:1033] map[call:craft\helpers\Db::batchInsert('{{%entries_authors}}', array(3), array(1)) file:/var/www/html/vendor/craftcms/cms/src/elements/Entry.php line:2424] map[call:craft\elements\Entry::_saveAuthors() file:/var/www/html/vendor/craftcms/cms/src/elements/Entry.php line:2328] map[call:craft\elements\Entry::afterSave(true) file:/var/www/html/vendor/craftcms/cms/src/services/Elements.php line:3682] map[call:craft\services\Elements::craft\services\{closure}() file:/var/www/html/vendor/craftcms/cms/src/services/Elements.php line:1170] map[call:craft\services\Elements::ensureBulkOp(instance of Closure) file:/var/www/html/vendor/craftcms/cms/src/services/Elements.php line:3498] map[call:craft\services\Elements::_saveElementInternal(instance of craft\elements\Entry, true, false, null, array(1), false, false, true) file:/var/www/html/vendor/craftcms/cms/src/services/Elements.php line:1251] map[call:craft\services\Elements::saveElement(instance of craft\elements\Entry) file:/var/www/html/vendor/craftcms/cms/src/gql/base/ElementMutationResolver.php line:173] map[call:craft\gql\base\ElementMutationResolver::saveElement(instance of craft\elements\Entry) file:/var/www/html/vendor/craftcms/cms/src/gql/resolvers/mutations/Entry.php line:97] map[call:craft\gql\resolvers\mutations\Entry::saveEntry(null, array(2), array(2), instance of GraphQL\Type\Definition\ResolveInfo) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:623] map[call:GraphQL\Executor\ReferenceExecutor::resolveFieldValueOrError(instance of GraphQL\Type\Definition\FieldDefinition, instance of GraphQL\Language\AST\FieldNode, array(2), null, instance of GraphQL\Type\Definition\ResolveInfo) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:549] map[call:GraphQL\Executor\ReferenceExecutor::resolveField(GraphQLType: Mutation, null, instance of ArrayObject(1), array(1)) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:474] map[call:GraphQL\Executor\ReferenceExecutor::GraphQL\Executor\{closure}(array(0), 'save_productCategories_productCategory_Entry') file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:857] map[call:GraphQL\Executor\ReferenceExecutor::GraphQL\Executor\{closure}(array(0), 'save_productCategories_productCategory_Entry')] map[file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php function:array_reduce(array(1), instance of Closure, array(0)) line:847] map[call:GraphQL\Executor\ReferenceExecutor::promiseReduce(array(1), instance of Closure, array(0)) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:468] map[call:GraphQL\Executor\ReferenceExecutor::executeFieldsSerially(GraphQLType: Mutation, null, array(0), instance of ArrayObject(1)) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:263] map[call:GraphQL\Executor\ReferenceExecutor::executeOperation(instance of GraphQL\Language\AST\OperationDefinitionNode, null) file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/ReferenceExecutor.php line:215] map[call:GraphQL\Executor\ReferenceExecutor::doExecute() file:/var/www/html/vendor/webonyx/graphql-php/src/Executor/Executor.php line:156] map[call:GraphQL\Executor\Executor::promiseToExecute(instance of GraphQL\Executor\Promise\Adapter\SyncPromiseAdapter, instance of GraphQL\Type\Schema, instance of GraphQL\Language\AST\DocumentNode, null, array(2), array(2), null, null) file:/var/www/html/vendor/webonyx/graphql-php/src/GraphQL.php line:161] map[call:GraphQL\GraphQL::promiseToExecute(instance of GraphQL\Executor\Promise\Adapter\SyncPromiseAdapter, instance of GraphQL\Type\Schema, '
                mutation InsertCategory($authorId: ID!, $title: String!) {
                        save_productCategories_productCategory_Entry(
                                authorId: $authorId
                                title: $title
                        ) {
                                authorId
                                title
                        }
                }
        ', null, array(2), array(2), null, null, array(26)) file:/var/www/html/vendor/webonyx/graphql-php/src/GraphQL.php line:93] map[call:GraphQL\GraphQL::executeQuery(instance of GraphQL\Type\Schema, '
                mutation InsertCategory($authorId: ID!, $title: String!) {
                        save_productCategories_productCategory_Entry(
                                authorId: $authorId
                                title: $title
                        ) {
                                authorId
                                title
                        }
                }
        ', null, array(2), array(2), null, null, array(26)) file:/var/www/html/vendor/craftcms/cms/src/services/Gql.php line:525] map[call:craft\services\Gql::executeQuery(instance of craft\models\GqlSchema, '
                mutation InsertCategory($authorId: ID!, $title: String!) {
                        save_productCategories_productCategory_Entry(
                                authorId: $authorId
                                title: $title
                        ) {
                                authorId
                                title
                        }
                }
        ', array(2), null, true) file:/var/www/html/vendor/craftcms/cms/src/controllers/GraphqlController.php line:194] map[call:craft\controllers\GraphqlController::actionApi()] map[file:/var/www/html/vendor/yiisoft/yii2/base/InlineAction.php function:call_user_func_array(array(2), array(0)) line:57] map[call:yii\base\InlineAction::runWithParams(array(0)) file:/var/www/html/vendor/yiisoft/yii2/base/Controller.php line:178] map[call:yii\base\Controller::runAction('api', array(0)) file:/var/www/html/vendor/yiisoft/yii2/base/Module.php line:552] map[call:yii\base\Module::runAction('graphql/api', array(0)) file:/var/www/html/vendor/craftcms/cms/src/web/Application.php line:349] map[call:craft\web\Application::runAction('graphql/api', array(0)) file:/var/www/html/vendor/yiisoft/yii2/web/Application.php line:103] map[call:yii\web\Application::handleRequest(instance of craft\web\Request) file:/var/www/html/vendor/craftcms/cms/src/web/Application.php line:317] map[call:craft\web\Application::handleRequest(instance of craft\web\Request) file:/var/www/html/vendor/yiisoft/yii2/base/Application.php line:384] map[call:yii\base\Application::run() file:/var/www/html/web/index.php line:12]]]]

```

- [ ] WIP: REST Testing in comparison - is this purely a GraphQL problem with mutations

```bash
# needs go installed on your local laptop
cd test-concurrent
go run concurrent.go -testType=rest   
```

Local DDEV testing: 
- sometimes "read tcp -> connection reset by peer"
- after first attempts: runs in `net/http: TLS handshake timeout` -> then needs restart of Orbstack (and restart DDEV project with `ddev start`), Docker network related I guess
- some (rare) actions were not successful with, but this was an exception

```
Raw Response Body: Bad Gateway
REST Request Error: error unmarshalling response body: invalid character 'B' looking for beginning of value
Raw Response Body: Internal Server Error
REST Request Error: error unmarshalling response body: invalid character 'I' looking for beginning of value
```

BUT ... no deadlock seen yet so far like in GraphQL testing. 

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