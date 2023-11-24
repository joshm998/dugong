import Theme from "./theme";
// import { MarkdownShortcutPlugin } from "@lexical/react/LexicalMarkdownShortcutPlugin";
import { ClientOnly } from 'remix-utils/client-only'
import { ContentEditable, LexicalComposer, RichTextPlugin, ErrorBoundary, TRANSFORMERS, MarkdownShortcutPlugin, HeadingNode, ListNode, ListItemNode, QuoteNode, CodeNode, CodeHighlightNode, LinkNode, AutoLinkNode, TableCellNode, TableNode, TableRowNode } from './lexical.client'
import ToolbarPlugin from "./toolbar";
// import ListMaxIndentLevelPlugin from "./richtext-plugins/max-indent";
// import CodeHighlightPlugin from "./richtext-plugins/code-highlight";
// import AutoLinkPlugin from "./richtext-plugins/autolink";

function Placeholder() {
  return <div className="editor-placeholder"></div>;
}

const editorConfig = {
  // The editor theme
  theme: Theme,
  // Handling of errors during update
  onError(error: any) {
    throw error;
  },
  namespace: "aaa",
  // Any custom nodes go here
  nodes: [
    HeadingNode,
    ListNode,
    ListItemNode,
    QuoteNode,
    CodeNode,
    CodeHighlightNode,
    TableNode,
    TableCellNode,
    TableRowNode,
    AutoLinkNode,
    LinkNode
  ]
};

const Editor = () => {
  return (
    <ClientOnly fallback={<p>Loading</p>}>
      {() => (
        <LexicalComposer initialConfig={editorConfig}>
          <div className="editor-container border rounded-lg min-h-[200px]">
            <ToolbarPlugin />
            <div className="editor-inner min-h-[200px]">
              <RichTextPlugin
                contentEditable={<ContentEditable className="editor-input p-2 min-h-[200px]" />}
                placeholder={<Placeholder />}
                ErrorBoundary={ErrorBoundary}
              />
              {/* <HistoryPlugin /> */}
              {/* <AutoFocusPlugin /> */}
              {/* <CodeHighlightPlugin /> */}
              {/* <ListPlugin /> */}
              {/* <LinkPlugin /> */}
              {/* <AutoLinkPlugin /> */}
              {/* <ListMaxIndentLevelPlugin maxDepth={7} /> */}
              <MarkdownShortcutPlugin transformers={TRANSFORMERS} />
            </div>
          </div>
        </LexicalComposer>
      )}
    </ClientOnly>
  );
}

export default Editor;