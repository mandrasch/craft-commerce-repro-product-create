<?php

namespace modules;

use Craft;
use craft\base\Event;
use craft\web\UrlManager;
use yii\base\Module as BaseModule;
use craft\events\RegisterUrlRulesEvent;

// TODO: Move everythin in folder my-module or myModule?

/**
 * MyModule module
 *
 * @method static MyModule getInstance()
 */
class MyModule extends BaseModule
{
    public function init(): void
    {
        Craft::setAlias('@modules', __DIR__);

        // TODO: is this needed for rest api routing?
        Craft::setAlias('@modules/my-module', $this->getBasePath());

        // Set the controllerNamespace based on whether this is a console or web request
        if (Craft::$app->request->isConsoleRequest) {
            $this->controllerNamespace = 'modules\\console\\controllers';
        } else {
            $this->controllerNamespace = 'modules\\controllers';
        }

        parent::init();

        $this->attachEventHandlers();

        // Any code that creates an element query or loads Twig should be deferred until
        // after Craft is fully initialized, to avoid conflicts with other plugins/modules
        Craft::$app->onInit(function() {
            // ...
        });
    }

    private function attachEventHandlers(): void
    {
        // Register event handlers here ...
        // (see https://craftcms.com/docs/5.x/extend/events.html to get started)

        // REST API
        // Beware - folders are like in the file system ("restApi"), but controller name needs to be in
        // dash-case ("project-pages" for ProjectPagesController) as well as the corresponding function
        // ("get-list" for "actionGetList")
        Event::on(UrlManager::class, UrlManager::EVENT_REGISTER_SITE_URL_RULES,
            function(RegisterUrlRulesEvent $event) {
                $event->rules = array_merge($event->rules, [
                    // https://craft-commerce-repro-product-create.ddev.site/rest-api/product-categories
                    'GET rest-api/product-categories' => 'my-module/restApi/product-categories/list',
                    'POST rest-api/product-categories/create' => 'my-module/restApi/product-categories/create',
                ]);
            }
        );
    }
}
