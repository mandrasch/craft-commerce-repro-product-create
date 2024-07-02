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

        $elementService = Craft::$app->getElements();

        $product = new Product();
        //  or use $product = $elementService->createElement(ProductElement::class); ?

        // use validated values from our model, set them to the entry we're about to save
        $product->title = 'My test product - '.time();
        $productTypeId = Commerce::getInstance()->getProductTypes()->getProductTypeByHandle('general')->id;
        $product->typeId = $productTypeId;

        // A simple variant
        $variant = new Variant();
        $variant->sku = time();
        $variant->price = 10;

        $product->setVariants([$variant]);

        // Set custom field (required field, we shouldn't get past save?)
        /*$product->setFieldValues([
            'customRequiredField1' => 'Test 123',
        ]);*/

        try {
            // Try to save new entry (with validation)
            if (!$elementService->saveElement($product, true)) {
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

        return ExitCode::OK;
    }
}