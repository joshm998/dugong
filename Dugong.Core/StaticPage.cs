namespace StaticSiteGen;

using RazorLight;

public abstract class StaticPage<TModel> : TemplatePage<TModel>
{
    public string CustomText { get; } = 
        "Gardyloo! - A Scottish warning yelled from a window before dumping" +
        "a slop bucket on the street below.";
}