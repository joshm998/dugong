namespace StaticSiteGen;

public class SiteModel
{
    public string SiteTitle { get; set; }
    public string SiteDescription { get; set; }
    
    public List<Content> AllContent { get; set; }
    
    public Content? Content { get; set; }
}