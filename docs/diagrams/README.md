# Diagram

De renderade PNG-bilderna i den här katalogen visas i [README.md](../../README.md). Källkoden (Mermaid) finns dels som HTML-kommentarer direkt ovanför respektive bild i README.md, dels som fristående `.mmd`-filer i [src/](src/) för enkel återrendering.

PNG:erna används i stället för att bädda in Mermaid-kodblock direkt, eftersom inte alla markdown-renderare (t.ex. på Bitbucket/interna wikis) stödjer Mermaid.

## Rendera om en bild

```sh
npx --yes -p @mermaid-js/mermaid-cli mmdc -i src/<namn>.mmd -o <namn>.png -b white -s 3
```

Uppdatera sedan samma källkod i motsvarande HTML-kommentar i README.md så att de två hålls i synk.
