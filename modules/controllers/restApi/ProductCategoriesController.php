<?php

namespace modules\controllers\restApi;

use Craft;
use craft\web\Controller;
use craft\elements\Entry;



class ProductCategoriesController extends Controller
{

    // allow anon
    protected int|bool|array $allowAnonymous = true;
    // allow methods
    protected array $verbs = ['POST', 'GET'];
    // TODO: restrict REST API via token with before() for real projects!
    // allow POST without CSRF
    public $enableCsrfValidation = false;

    public function actionList(): \yii\web\Response
    {
        // mock data
        return $this->asJson([
            'success' => true,
            'data' => [
                [
                    'id' => 1,
                    'name' => 'Category 1',
                ],
                [
                    'id' => 2,
                    'name' => 'Category 2',
                ],
            ],
        ]);
    }

    public function actionCreate(): \yii\web\Response
    {
        // TODO: require POST
        $elementService = Craft::$app->getElements();
        $entry = $elementService->createElement(Entry::class);

        $section = Craft::$app->entries->getSectionByHandle('productCategories');
        $entryType = $section->getEntryTypes()[0];

        $entry->sectionId = $section->id;
        $entry->typeId = $entryType->id;

        $entry->title = $this->request->getBodyParam('title');

        // execute save, validated model values will be saved to entry ($product)
        if(!Craft::$app->elements->saveElement($entry)){
            // return 400 as json
            return $this->asJson([
                'success' => false,
                'errors' => $entry->getErrors(),
            ]);
        }
        return $this->asJson([
            'success' => true,
        ]);
    }

}
