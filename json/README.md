# JSON Data Structure: Brizdoors

This repo's `json/` folder is organized for a robust, extensible product backend, ready for both developer and AI/automation usage.

## Structure

- products/
  - schemas/                # core JSON schema/templates only, not real data
    - hardware_schema.json  # hardware data definition
    - door_schema.json      # door product data definition
    - ...

  - hardware/               # each file is a real hardware product, named by unique SKU/code
    - HW-LEGGE-ENT-001.json
    - ...

  - finishes/               # all real finish codes, one file: finishes.json
    - finishes.json

  - doors/                  # door product definitions, per file
    - DOOR-SCMDF-LEG01-ENT.json
    - ...

  - lockwood_assa_abloy/    # per-manufacturer/brand subfolders allowed (optional, for large catalogs)
      1220_series/
        1220.json
      ...

## Linking

- All product/hardware/door/finish relations must use an explicit unique code/key as defined in the relevant schema.
- Enums, types, and allowable values should be consistent and documented in README/schema.

## Example Workflows

A product references its finish by finish_code and hardware by hardware_item_code; the actual data objects are in `finishes.json` and `/hardware/*.json`.

## Maintenance
- ALL template/sample schemas go in /products/schemas, not mixed with real data.
- If you add a new product type or manufacturer, add a brief README there.
- Use one JSON object per file for items; lists only in code tables (like all finishes).

---

