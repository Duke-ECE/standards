# Migrations for the shared Supabase project

Several repos push migrations to the same project, so version files by
**timestamp**, never plain counters:

```
supabase/migrations/20260726000001_create_widgets.sql
```

Rules: schema only (no secrets, no seed data with keys), `create table if not
exists`, `enable row level security` with no policies for service-owned
tables. Apply with `supabase link --project-ref <ref>` + `supabase db push`.
