# Official Apple Ads API Docs (Index + CLI Mapping)

This repo implements a CLI (`aads`) for Apple's **Apple Ads Campaign Management API 5**.

This document intentionally **links to** official Apple Developer Documentation instead of copying it verbatim.

For a complete, link-only index of *all* Apple Ads docs pages (endpoints + objects), see: `docs/OFFICIAL_APPLE_ADS_DOC_INDEX.md`.

## Essentials

- OAuth: [Implementing OAuth for the Apple Ads API](https://developer.apple.com/documentation/apple_ads/implementing-oauth-for-the-apple-search-ads-api)
- Request headers + rate limits: [Calling the Apple Ads API](https://developer.apple.com/documentation/apple_ads/calling-the-apple-search-ads-api)
- Selectors + partial fetch: [Using Apple Ads API Functionality](https://developer.apple.com/documentation/apple_ads/using-apple-search-ads-api-functionality)

## Endpoint Mapping

### Apps

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads apps search` | `GET /search/apps` | [Search for iOS apps](https://developer.apple.com/documentation/apple_ads/search-for-ios-apps) |
| `aads apps eligibility` | `POST /app-eligibility/find` | [Find App Eligibility Records](https://developer.apple.com/documentation/apple_ads/find-app-eligibility-records) |
| `aads apps details` | `GET /apps/{adamId}` | [Get App Details](https://developer.apple.com/documentation/apple_ads/get-app-details) |
| `aads apps localized` | `GET /apps/{adamId}/localized-details` | [Get Localized App Details](https://developer.apple.com/documentation/apple_ads/get-localized-app-details) |

### Campaigns

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads campaigns create` | `POST /campaigns` | [Create a Campaign](https://developer.apple.com/documentation/apple_ads/create-a-campaign) |
| `aads campaigns find` | `POST /campaigns/find` | [Find Campaigns](https://developer.apple.com/documentation/apple_ads/find-campaigns) |
| `aads campaigns get` | `GET /campaigns/{campaignId}` | [Get a Campaign](https://developer.apple.com/documentation/apple_ads/get-a-campaign) |
| `aads campaigns list` | `GET /campaigns` | [Get all Campaigns](https://developer.apple.com/documentation/apple_ads/get-all-campaigns) |
| `aads campaigns update` | `PUT /campaigns/{campaignId}` | [Update a Campaign](https://developer.apple.com/documentation/apple_ads/update-a-campaign) |
| `aads campaigns delete` | `DELETE /campaigns/{campaignId}` | [Delete a Campaign](https://developer.apple.com/documentation/apple_ads/delete-a-campaign) |

### Budget Orders

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads budgetorders create` | `POST /budgetorders` | [Create a Budget Order](https://developer.apple.com/documentation/apple_ads/create-a-budget-order) |
| `aads budgetorders update` | `PUT /budgetorders/{budgetOrderId}` | [Update a Budget Order](https://developer.apple.com/documentation/apple_ads/update-a-budget-order) |
| `aads budgetorders get` | `GET /budgetorders/{budgetOrderId}` | [Get a Budget Order](https://developer.apple.com/documentation/apple_ads/get-a-budget-order) |
| `aads budgetorders list` | `GET /budgetorders` | [Get All Budget Orders](https://developer.apple.com/documentation/apple_ads/get-all-budget-orders) |

### Ad Groups

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads adgroups create` | `POST /campaigns/{campaignId}/adgroups` | [Create an Ad Group](https://developer.apple.com/documentation/apple_ads/create-an-ad-group) |
| `aads adgroups find` | `POST /campaigns/{campaignId}/adgroups/find` | [Find Ad Groups](https://developer.apple.com/documentation/apple_ads/find-ad-groups) |
| `aads adgroups find-all` | `POST /adgroups/find` | [Find Ad Groups (org-level)](https://developer.apple.com/documentation/apple_ads/find-ad-groups-(org-level)) |
| `aads adgroups get` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}` | [Get an Ad Group](https://developer.apple.com/documentation/apple_ads/get-an-ad-group) |
| `aads adgroups list` | `GET /campaigns/{campaignId}/adgroups` | [Get all Ad Groups](https://developer.apple.com/documentation/apple_ads/get-all-ad-groups) |
| `aads adgroups update` | `PUT /campaigns/{campaignId}/adgroups/{adgroupId}` | [Update an Ad Group](https://developer.apple.com/documentation/apple_ads/update-an-ad-group) |
| `aads adgroups delete` | `DELETE /campaigns/{campaignId}/adgroups/{adgroupId}` | [Delete an Ad Group](https://developer.apple.com/documentation/apple_ads/delete-an-ad-group) |

### Targeting Keywords

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads keywords create` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords/bulk` | [Create Targeting Keywords](https://developer.apple.com/documentation/apple_ads/create-targeting-keywords) |
| `aads keywords find-campaign` | `POST /campaigns/{campaignId}/targetingkeywords/find` | [Find Targeting Keywords in a Campaign](https://developer.apple.com/documentation/apple_ads/find-targeting-keywords-in-a-campaign) |
| `aads keywords get` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords/{keywordId}` | [Get a Targeting Keyword in an Ad Group](https://developer.apple.com/documentation/apple_ads/get-a-targeting-keyword-in-an-ad-group) |
| `aads keywords list` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords` | [Get All Targeting Keywords in an Ad Group](https://developer.apple.com/documentation/apple_ads/get-all-targeting-keywords-in-an-ad-group) |
| `aads keywords update` | `PUT /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords/bulk` | [Update Targeting Keywords](https://developer.apple.com/documentation/apple_ads/update-targeting-keywords) |
| `aads keywords delete` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords/delete/bulk` | [Delete Targeting Keywords](https://developer.apple.com/documentation/apple_ads/delete-targeting-keywords) |
| `aads keywords delete-one` | `DELETE /campaigns/{campaignId}/adgroups/{adgroupId}/targetingkeywords/{keywordId}` | [Delete a Targeting Keyword](https://developer.apple.com/documentation/apple_ads/delete-a-targeting-keyword) |

### Negative Keywords (Campaign-Level)

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads negatives campaign-create` | `POST /campaigns/{campaignId}/negativekeywords/bulk` | [Create Campaign Negative Keywords](https://developer.apple.com/documentation/apple_ads/create-campaign-negative-keywords) |
| `aads negatives campaign-find` | `POST /campaigns/{campaignId}/negativekeywords/find` | [Find Campaign Negative Keywords](https://developer.apple.com/documentation/apple_ads/find-campaign-negative-keywords) |
| `aads negatives campaign-get` | `GET /campaigns/{campaignId}/negativekeywords/{keywordId}` | [Get a Campaign Negative Keyword](https://developer.apple.com/documentation/apple_ads/get-a-campaign-negative-keyword) |
| `aads negatives campaign-list` | `GET /campaigns/{campaignId}/negativekeywords` | [Get All Campaign Negative Keywords](https://developer.apple.com/documentation/apple_ads/get-all-campaign-negative-keywords) |
| `aads negatives campaign-update` | `PUT /campaigns/{campaignId}/negativekeywords/bulk` | [Update Campaign Negative Keywords](https://developer.apple.com/documentation/apple_ads/update-campaign-negative-keywords) |
| `aads negatives campaign-delete` | `POST /campaigns/{campaignId}/negativekeywords/delete/bulk` | [Delete Campaign Negative Keywords](https://developer.apple.com/documentation/apple_ads/delete-campaign-negative-keywords) |

### Negative Keywords (Ad Group-Level)

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads negatives adgroup-create` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords/bulk` | [Create Ad Group Negative Keywords](https://developer.apple.com/documentation/apple_ads/create-ad-group-negative-keywords) |
| `aads negatives adgroup-find` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords/find` | [Find Ad Group Negative Keywords](https://developer.apple.com/documentation/apple_ads/find-ad-group-negative-keywords) |
| `aads negatives adgroup-get` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords/{keywordId}` | [Get an Ad Group Negative Keyword](https://developer.apple.com/documentation/apple_ads/get-an-ad-group-negative-keyword) |
| `aads negatives adgroup-list` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords` | [Get All Ad Group Negative Keywords](https://developer.apple.com/documentation/apple_ads/get-all-ad-group-negative-keywords) |
| `aads negatives adgroup-update` | `PUT /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords/bulk` | [Update Ad Group Negative Keywords](https://developer.apple.com/documentation/apple_ads/update-ad-group-negative-keywords) |
| `aads negatives adgroup-delete` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/negativekeywords/delete/bulk` | [Delete Ad Group Negative Keywords](https://developer.apple.com/documentation/apple_ads/delete-ad-group-negative-keywords) |

### Ads

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads ads create` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/ads` | [Create an Ad](https://developer.apple.com/documentation/apple_ads/create-an-ad) |
| `aads ads find` | `POST /campaigns/{campaignId}/adgroups/{adgroupId}/ads/find` | [Find Ads](https://developer.apple.com/documentation/apple_ads/find-ads) |
| `aads ads find-all` | `POST /ads/find` | [Find Ads (org-level)](https://developer.apple.com/documentation/apple_ads/find-ads-(org-level)) |
| `aads ads get` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/ads/{adId}` | [Get an Ad](https://developer.apple.com/documentation/apple_ads/get-an-ad) |
| `aads ads list` | `GET /campaigns/{campaignId}/adgroups/{adgroupId}/ads` | [Get All Ads](https://developer.apple.com/documentation/apple_ads/get-all-ads) |
| `aads ads update` | `PUT /campaigns/{campaignId}/adgroups/{adgroupId}/ads/{adId}` | [Update an Ad](https://developer.apple.com/documentation/apple_ads/update-an-ad) |
| `aads ads delete` | `DELETE /campaigns/{campaignId}/adgroups/{adgroupId}/ads/{adId}` | [Delete an Ad](https://developer.apple.com/documentation/apple_ads/delete-an-ad) |

### Creatives

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads creatives create` | `POST /creatives` | [Create a Creative](https://developer.apple.com/documentation/apple_ads/create-a-creative) |
| `aads creatives find` | `POST /creatives/find` | [Find Creatives](https://developer.apple.com/documentation/apple_ads/find-creatives) |
| `aads creatives get` | `GET /creatives/{creativeId}` | [Get a Creative](https://developer.apple.com/documentation/apple_ads/get-a-creative) |
| `aads creatives list` | `GET /creatives` | [Get All Creatives](https://developer.apple.com/documentation/apple_ads/get-all-creatives) |

### Custom Product Pages

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads product-pages list` | `GET /apps/{adamId}/product-pages` | [Get Product Pages](https://developer.apple.com/documentation/apple_ads/get-product-pages) |
| `aads product-pages get` | `GET /apps/{adamId}/product-pages/{productPageId}` | [Get Product Pages by Identifier](https://developer.apple.com/documentation/apple_ads/get-product-pages-by-identifier) |
| `aads product-pages locales` | `GET /apps/{adamId}/product-pages/{productPageId}/locale-details` | [Get Product Page Locales](https://developer.apple.com/documentation/apple_ads/get-product-page-locales) |
| `aads product-pages countries` | `GET /countries-or-regions` | [Get Supported Countries or Regions](https://developer.apple.com/documentation/apple_ads/get-supported-countries-or-regions) |
| `aads product-pages device-sizes` | `GET /creativeappmappings/devices` | [Get App Preview Device Sizes](https://developer.apple.com/documentation/apple_ads/get-app-preview-device-sizes) |

### Ad Rejection Reasons

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads ad-rejections find` | `POST /product-page-reasons/find` | [Find Ad Creative Rejection Reasons](https://developer.apple.com/documentation/apple_ads/find-ad-creative-rejection-reasons) |
| `aads ad-rejections get` | `GET /product-page-reasons/{productPageReasonId}` | [Gets a Product Page Reason](https://developer.apple.com/documentation/apple_ads/gets-a-product-page-reason) |
| `aads ad-rejections find-assets` | `POST /apps/{adamId}/assets/find` | [Find App Assets](https://developer.apple.com/documentation/apple_ads/find-app-assets) |

### Search Geolocations

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads geo search` | `GET /search/geo` | [Search for Geolocations](https://developer.apple.com/documentation/apple_ads/search-for-geolocations) |
| `aads geo get` | `GET /geodata` | [Get a List of Geo Locations](https://developer.apple.com/documentation/apple_ads/get-a-list-of-geo-locations) |

### Reports

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads reports campaigns` | `POST /reports/campaigns` | [Get Campaign-Level Reports](https://developer.apple.com/documentation/apple_ads/get-campaign-level-reports) |
| `aads reports adgroups` | `POST /reports/campaigns/{campaignId}/adgroups` | [Get Ad Group-Level Reports](https://developer.apple.com/documentation/apple_ads/get-ad-group-level-reports) |
| `aads reports keywords` | `POST /reports/campaigns/{campaignId}/keywords` | [Get Keyword-Level Reports](https://developer.apple.com/documentation/apple_ads/get-keyword-level-reports) |
| `aads reports keywords --adgroup-id ...` | `POST /reports/campaigns/{campaignId}/adgroups/{adgroupId}/keywords` | [Get Keyword-Level Within Ad Group Reports](https://developer.apple.com/documentation/apple_ads/get-keyword-level-within-ad-group-reports) |
| `aads reports searchterms` | `POST /reports/campaigns/{campaignId}/searchterms` | [Get Search Term-Level Reports](https://developer.apple.com/documentation/apple_ads/get-search-term-level-reports) |
| `aads reports searchterms --adgroup-id ...` | `POST /reports/campaigns/{campaignId}/adgroups/{adgroupId}/searchterms` | [Get Search Term-Level Within Ad Group Reports](https://developer.apple.com/documentation/apple_ads/get-search-term-level-within-ad-group-reports) |
| `aads reports ads` | `POST /reports/campaigns/{campaignId}/ads` | [Get Ad-Level Reports](https://developer.apple.com/documentation/apple_ads/get-ad-level-reports) |

### Impression Share Reports

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads impression-share create` | `POST /custom-reports` | [Impression Share Report](https://developer.apple.com/documentation/apple_ads/impression-share-report) |
| `aads impression-share get` | `GET /custom-reports/{customReportId}` | [Get a Single Impression Share Report](https://developer.apple.com/documentation/apple_ads/get-a-single-impression-share-report) |
| `aads impression-share list` | `GET /custom-reports` | [Get All Impression Share Reports](https://developer.apple.com/documentation/apple_ads/get-all-impression-share-reports) |

### ACLs

| CLI | HTTP | Apple Docs |
|---|---|---|
| `aads acls list` | `GET /acls` | [Get User ACL](https://developer.apple.com/documentation/apple_ads/get-user-acl) |
| `aads acls me` | `GET /me` | [Get Me Details](https://developer.apple.com/documentation/apple_ads/get-me-details) |
