# Översikt

Välkommen till den inbyggda hjälpen där förhoppningen är att du ska hitta nödvändig information för att hantera dina enheter och segment. I menyn till vänster hittar du länkar till olika hjälpavsnitt.

## Har du använt verktyget tidigare?

Om du är van användare av detta verktyg kan du hoppa direkt till hantering av enheter och system:

- [Lägg till ny endpoint](03-new-endpoint.md)
- [Lägg till ny hub](02-new-hub.md))
- [Lägg till nytt system](01-new-system.md)
- [Network Configuration](networking/overview.md)
- [Troubleshooting](troubleshooting/common-issues.md)

## Är du ny användare?

Som ny användare bör du läsa igenom följande avsnitt för att bekanta dig med verktygets upplägg och funktioner.

Verktygets huvudfunktion är att lägga på ett skyddande lager på befintliga, platta nätverk med en en central hub och flera endpoints. Inom mindre OT-system är det inte ovanligt med ett eller fler platta och oskyddade nätverk där ett eller flera är utspritt geografiskt med en tredje part som leverantör av kommunikationsinfrastrukturen. Det finns ett par sårbarheter med sådana nätverk där åtkomst vid en fysisk plats ger nätverksåtkomst till alla andra platser, både centralt och i fält. En annan sårbarhet är möjligheten att påverka den oskyddade kommunikationen via tredje parts kommunikationsinfrastruktur.

Genom att tillföra Teltonika routrar vid varje anslutningspunkt i fält och centralt skapar verktyget ett skyddat stjärnformat nätverk genom vilket den vanliga trafiken går, utan att påverka befintlig adressering. Vid den centrala enheten filtreras trafiken så att endast trafik som skall gå till respektive enhet i fält tillåts, övrig trafik blockeras.

![Typiskt kontrollernät](images/SCADA%20flat%20network-1.png)

### Begrepp

- **System**: Ett system i detta verktyg är inte mer än en samling enheter som tillsammans utgör ett skyddat nätverk.

- [Ny enhet](getting-started/overview.md)
- [Device Dashboard](getting-started/dashboard.md)
- [Device Management](devices/overview.md)
- [Network Configuration](networking/overview.md)
- [Troubleshooting](troubleshooting/common-issues.md)

## Vad är CDF Console?

CDF Console är tänkt att hjälpa er sätta upp krypterade mikrosegment på ett befintligt, platt närverk utan att påverka befintliga adressering genom att tillföra enklare men kapabla enheter med hjälpa av VPN (Wireguard) och filter (ebtables) för att dels hindra angrepp via kommunikationsinfrastrukturen samt angrepp av alla fältenheter på nätverk från en plats i fält.

## Att använda hjälpen

Menyn till vänster ger en översikt över olika hjälpavsnitt, men det finns även hjälpknappar i specifika vyer i verktyget som leder till motsvarande hjälpavsnitt.
