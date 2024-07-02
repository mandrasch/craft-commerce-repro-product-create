<?php

namespace modules\console\controllers;

use Craft;
use craft\console\Controller;
use craft\helpers\Console;
use yii\console\ExitCode;
use craft\commerce\elements\Product;
use craft\commerce\Plugin as Commerce;
use craft\commerce\elements\Variant;

/**
 * Greet Controller
 */
class ProductController extends Controller
{
    public $defaultAction = 'create';
    /**
     * Issues a greeting to new Craft developers.
     */
    public function actionCreate(): int
    {

        // Thanks to https://github.com/craftcms/commerce/discussions/3458

        // Prepare new entry
        $product = new Product();

        // Use validated values from our model, set them to the entry we're about to save
        $product->title = 'My test product - '.time();
        $productTypeId = Commerce::getInstance()->getProductTypes()->getProductTypeByHandle('general')->id;
        $product->typeId = $productTypeId;

        // Set custom field (required field, we shouldn't get past save?)
        /*$product->setFieldValues([
            'customRequiredField1' => 'Test 123',
        ]);*/

        try {
            // Try to save new entry (with validation = true by default)
            if (!Craft::$app->elements->saveElement($product)) {
                if ($product->hasErrors()) {
                    $validationErrors = [];

                    foreach ($product->getFirstErrors() as $attribute => $errorMessage) {
                        $validationErrors[] = $errorMessage;
                    }

                    $this->sdterr('Could not save product: '.implode(', ', $validationErrors), Console::FG_RED);
                 
                } else {
                    $this->sdterr('Could not save product', Console::FG_RED);
                }
            }
        } catch (\Throwable $e) {
            $this->sdterr('Could not save product', Console::FG_RED);
        }

        $this->stdout('Product created - id: '.$product->id.PHP_EOL, Console::FG_GREEN);

        // A simple variant
        $variant = new Variant();
        $variant->sku = time();
        $variant->price = 10;
        $variant->isDefault = true;
        $variant->title = 'My test variant - '.time();
        // important - otherwise error (Exception 'TypeError' with message 'craft\commerce\elements\Variant::updateTitle(): Argument #1 ($product) must be of type craft\commerce\elements\Product, null given, called in /var/www/html/vendor/craftcms/commerce/src/elements/Variant.php on line 1005')
        $variant->primaryOwnerId = $product->id;

        // Save variant before adding it to the product (important)
        if(!Craft::$app->elements->saveElement($variant)){
                $this->stderr('Could not save variant', Console::FG_RED);

                $this->stderr('Errors: '.print_r($variant->getErrors(), true));

                return ExitCode::UNSPECIFIED_ERROR;
        }
        $this->stdout('Variant saved'.PHP_EOL, Console::FG_GREEN);

        // Add variant to product
        $product->setVariants([$variant]);
        if(!Craft::$app->elements->saveElement($product)){
            $this->stderr('Could not save product', Console::FG_RED);
            return ExitCode::UNSPECIFIED_ERROR;
        }
        $this->stdout('Product saved'.PHP_EOL, Console::FG_GREEN);

        return ExitCode::OK;
    }
}