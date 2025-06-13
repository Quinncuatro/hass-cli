# Home Assistant Naming Strategy for hass-cli

This guide provides best practices for naming areas and entities in Home Assistant to maximize effectiveness with the hass-cli fuzzy search functionality.

## Area Naming Best Practices

**Use simple, single words when possible:**
```
✅ Good:
- office
- kitchen  
- bedroom
- living
- garage

❌ Avoid:
- henry's office
- living room area
- master bedroom suite
```

**For multi-word areas, use common abbreviations:**
```
✅ Good:
- master (for master bedroom)
- guest (for guest bedroom) 
- dining (for dining room)
```

## Entity Naming Patterns

**Follow this pattern: `[Area] [Type] [Location/Description]`**

### Lights
```
✅ Excellent:
- "Office Desk Lamp"
- "Office Ceiling"  
- "Kitchen Island"
- "Bedroom Nightstand"
- "Living Floor Lamp"

❌ Avoid:
- "Henry's Office - Desk - Floor Lamp"
- "Main light in the kitchen area"
- "Bedroom lamp (left side)"
```

### Other Devices
```
✅ Fans:
- "Office Ceiling Fan"
- "Bedroom Fan"

✅ Climate:
- "Office Thermostat"
- "Living AC"

✅ Switches:
- "Office Outlet"
- "Kitchen Coffee Maker"
```

## Entity ID Best Practices

**Match your friendly names:**
```
✅ Good:
friendly_name: "Office Desk Lamp"
entity_id: light.office_desk_lamp

friendly_name: "Kitchen Island"  
entity_id: light.kitchen_island
```

## Special Cases

**For areas with multiple similar devices, add descriptors:**
```
✅ Multiple lights in same room:
- "Office Desk Lamp"
- "Office Ceiling"
- "Office Bookshelf"

✅ Multiple bedrooms:
- "Master Desk Lamp"
- "Guest Ceiling"
```

## CLI Commands That Will Work Great

With this naming scheme, you'll get excellent results:

```bash
# These will work perfectly:
./bin/hass office lights on
./bin/hass kitchen lights off
./bin/hass bedroom fan 75
./bin/hass living temperature 72

# Multiple matches resolved intelligently:
./bin/hass office desk on      # Matches "Office Desk Lamp"
./bin/hass office ceiling off  # Matches "Office Ceiling"
```

## Home Assistant Area Assignment

**Pro tip:** Assign entities to areas in Home Assistant UI:
1. Go to Settings → Areas & Zones
2. Create areas: "Office", "Kitchen", "Bedroom", etc.
3. Assign devices to areas

*Note: The current CLI doesn't use area_id yet (since it's not in API responses), but this future-proofs your setup and helps Home Assistant's native features.*

## Quick Migration Strategy

**Priority order for renaming:**
1. **Lights first** (most commonly controlled)
2. **Fans and climate** (frequently used)
3. **Switches and sensors** (less critical)

**You don't need to rename everything at once!** The current fuzzy matching is quite good, so rename entities as you encounter matching issues.

## Technical Notes

The hass-cli fuzzy matching system:
- Scores entities based on domain matching (40%), area matching (30%), and name matching (30%)
- Gives higher scores to exact word matches in entity names
- Penalizes unavailable entities by 5%
- Adds small bonuses for more descriptive names (more words = more specific)
- Prefers entities that are currently available vs unavailable

Default fuzzy threshold is 0.6, but can be adjusted in `~/.config/hass/config.yaml`:
```yaml
preferences:
  fuzzy_threshold: 0.6  # Lower = more permissive matching
```