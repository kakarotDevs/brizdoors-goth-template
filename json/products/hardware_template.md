# Hardware Item JSON Template

## hardware_item_code
- Type: string
- Description: Unique internal SKU/ID for the item
- Example: `"LEG-G2-L52-PAS-SC"`

## product_type
- Type: string
- Description: Standard product type
- Example: `"lever_on_rose"`

## name
- Type: string
- Description: Human-readable name of the item
- Example: `"Legge G2 Passage Lever on Rose, Satin Chrome"`

## description
- Type: string
- Description: Detailed product description

## status
- Type: enum (`"active"`, `"obsolete"`, `"draft"`, `"discontinued"`, `"special_order"`)
- Default: `"active"`

---

## manufacturer_details
- name: string (e.g. `"Legge"`)
- part_code: string or null
- article_number: string or null
- series: string or null
- ordering_method: string (e.g. `"Distributor Portal"`)
- order_reference_code_format: string (e.g. `"Series/Lever/Function/Finish"`)

---

## finish_options
- Type: array of strings
- Example: `["SC", "CP", "PB"]`

## selected_finish_code
- Type: string or null

---

## backset_options
- Type: array of objects `{ "value": number, "unit": string }`
- Example: `[ { "value": 60, "unit": "mm" }, { "value": 70, "unit": "mm" } ]`

## selected_backset
- Type: object or null

---

## requires_handing
- Type: boolean or null

## handing
- Type: string (`"L"`, `"R"`, `"U"`) or null

---

## standardized_function_description
- Type: string

## is_as1428_compliant
- Type: boolean or null

## certifications
- Type: array of strings

## tags
- Type: array of strings

## image_url
- Type: string

## material_main
- Type: string or null

---

## lever_on_rose_details
- lever_style_code: string
- rose_style_code: string
- manufacturer_rose_no: string or null
- rose_type: string or null (e.g. `"interior"`)

---

## plate_furniture_details
- plate_series_or_style_code: string
- plate_function_code_manufacturer: string
- handle_style_code: string

---

## mortice_lock_details
- manufacturer_function_code: string
- manufacturer_lock_function_term: string
- compatible_door_material_code_manufacturer: string
- intrinsic_backset: object or null
- case_dimensions_mm: object with `height`, `depth`, `thickness`
- faceplate_dimensions_mm: object with `length`, `width`
- strike_plate_included: boolean or null
- spindle_size_mm: number or null
- cylinder_type_compatible: string

---

## ordering_info
- requires_combination_with_other_parts: boolean or null
- example_order_string_manufacturer: string
- notes_for_ordering: string

---

## compatibility
- compatible_door_series_internal: array of strings
- compatible_hardware_item_codes_internal: array of strings

---

## additional_attributes_raw
- original_finish_term_manufacturer: string
- original_function_term_manufacturer: string
- brochure_page_no: number or null
- manufacturer_specific_notes: string

---

## dimensions_overall_package
- width, height, depth, weight: objects with `value` and `unit`

---

## visibility_scope
- show_on_client_portal: boolean
- available_for_quote_generation: boolean
- internal_use_only: boolean

---

## metadata
- last_updated: string (timestamp) or null
- created_date: string (timestamp) or null
- version: number
- internal_notes: string
