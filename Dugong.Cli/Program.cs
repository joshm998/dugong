using Spectre.Console;
using Spectre.Console.Cli;
using StaticSiteGen.Cli.DependencyInjection;

var services = new ServiceCollection();

services.AddScoped<BuildService>();

var app = new CommandApp(new TypeRegistrar(services));

app.Configure(c =>
{
    c.AddCommand<BuildSite>("build");
    c.AddCommand<BuildSite>("watch");
});

if (args.Length == 0)
{
    AnsiConsole.Write(new FigletText("StaticSiteGen").Centered().Color(Color.Purple));
}

await app.RunAsync(args);