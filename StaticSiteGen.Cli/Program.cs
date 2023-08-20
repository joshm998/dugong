using System.Net;
using System.Text.RegularExpressions;
using Markdig;
using Markdig.Extensions.Yaml;
using Markdig.Renderers;
using Markdig.Syntax;
using RazorLight;
using StaticSiteGen;

var path = Directory.GetCurrentDirectory();

path = Directory.GetParent(path)?.FullName;

var directory = $"{path}/StaticSiteGen.DocsSite";

var engine = new RazorLightEngineBuilder()
    .UseFileSystemProject(directory)
    .UseMemoryCachingProvider()
    .Build();

var parsedContent = new List<Content>();

var content = Directory.GetFiles($"{directory}/content", "*.md", SearchOption.AllDirectories);

foreach (var file in content)
{
    var relativePath = file
        .Replace($"{directory}/content/", "");

    var markdown = File.ReadAllText(file);
    var pipeline = new MarkdownPipelineBuilder()
        .UseYamlFrontMatter()
        .Build();
    
    StringWriter writer = new StringWriter();
    var renderer = new HtmlRenderer(writer);
    pipeline.Setup(renderer);
    
    MarkdownDocument document = Markdown.Parse(markdown, pipeline);

    var yamlBlock = document.Descendants<YamlFrontMatterBlock>().FirstOrDefault();
    
    if (yamlBlock != null)
    {
        string yaml = markdown.Substring(yamlBlock.Span.Start, yamlBlock.Span.Length);
    
        Console.WriteLine(yaml);
    }
    
    renderer.Render(document);
    writer.Flush();
    string html = writer.ToString();

    var a = new Content()
    {
        Slug = Path.GetFileNameWithoutExtension(file),
        ContentPath = relativePath.Replace($"/{Path.GetFileName(file)}", ""),
        Html = WebUtility.HtmlDecode(html)
    };
    
    parsedContent.Add(a);
}

//TODO: Fetch From JSON
var siteModel = new SiteModel() { SiteTitle = "StaticSiteGen", SiteDescription = "Test", AllContent = parsedContent };


var pages = Directory.GetFiles($"{directory}/pages", "*.cshtml", SearchOption.AllDirectories);

if (!Directory.Exists($"{directory}/output"))
{
    Directory.CreateDirectory($"{directory}/output");
}

foreach (var page in pages)
{

    var relativePath = page
        .Replace($"{directory}/pages/", "")
        .Replace(Path.GetExtension(page), "");

    var match = Regex.Match(relativePath, @"\[(.*?)\]");
    if (match.Success)
    {
        var contentPath = relativePath.Replace($"/{Path.GetFileNameWithoutExtension(page)}", "");
        var matchingContent = parsedContent
            .Where(e => e.ContentPath == contentPath);

        foreach (var a in matchingContent)
        {
            var nameSelector = match.Groups[1].Value;
            string? pageName = null;
            
            
            if (nameSelector == "Slug")
                pageName = a.Slug;
                
            if (pageName == null)
            {
                throw new Exception("No Selector Found... Please use a valid value.");
            }

            var pageModel = siteModel;
            pageModel.Content = a;
            
            var result = await engine.CompileRenderAsync(page.Replace($"{directory}/", ""), pageModel);
            
            var generatedPageName = $"{directory}/output/{contentPath}/{pageName}.html".ToLower();
            
            var file = new FileInfo(generatedPageName);
            file.Directory?.Create();
            await using (var sw = File.CreateText(generatedPageName))
            {
                await sw.WriteAsync(result);
            }
            Console.WriteLine($"Page Created - {generatedPageName}");
        }

    }
    else
    {
        var result = await engine.CompileRenderAsync(page.Replace($"{directory}/", ""), siteModel);
        var generatedPageName = $"{directory}/output/{relativePath.ToLower()}.html";
        var file = new FileInfo(generatedPageName);
        file.Directory?.Create();
        await using (var sw = File.CreateText(generatedPageName))
        {
            await sw.WriteAsync(result);
        }
        Console.WriteLine($"Page Created - {relativePath.ToLower()}");
    }
}