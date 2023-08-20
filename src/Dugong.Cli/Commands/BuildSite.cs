using RazorLight;
using Spectre.Console;
using Spectre.Console.Cli;

public class BuildSiteSettings : CommandSettings
{
}

public class BuildSite : AsyncCommand<BuildSiteSettings>
{
    private readonly BuildService _buildService;
    public BuildSite(BuildService buildService)
    {
        _buildService = buildService;
    }

    public override async Task<int> ExecuteAsync(CommandContext context, BuildSiteSettings settings)
    {
        Tree? fileGeneratedTree = null;

        await AnsiConsole.Status()
            .Spinner(Spinner.Known.Dots)
            .StartAsync("Building Site...", async ctx =>
            {
                var path = Directory.GetCurrentDirectory();

                path = Directory.GetParent(path)?.FullName;

                var directory = $"{path}/StaticSiteGen.DocsSite";

                var siteModel = await _buildService.GenerateContent(directory);

                var engine = new RazorLightEngineBuilder()
                    .UseFileSystemProject(directory)
                    .UseMemoryCachingProvider()
                    .Build();
                
                fileGeneratedTree = await _buildService.BuildTemplate(directory, engine, siteModel);
            });
        
        // Render the tree
        AnsiConsole.MarkupLine(":check_mark: Build: Successful");
        if (fileGeneratedTree != null)
            AnsiConsole.Write(fileGeneratedTree);
        return 0;
    }
}