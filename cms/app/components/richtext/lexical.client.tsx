// components/lexical.client.tsx
export * from 'lexical'
export * from '@lexical/react/LexicalComposer'
export * from '@lexical/react/LexicalRichTextPlugin'
export * from '@lexical/react/LexicalErrorBoundary'
export * from '@lexical/react/LexicalContentEditable'
export * from '@lexical/react/LexicalMarkdownShortcutPlugin'
export * from '@lexical/markdown'
export * from '@lexical/list'
export * from '@lexical/table'
export * from '@lexical/code'
export * from '@lexical/link'
export * from '@lexical/rich-text'
export {$isParentElementRTL, $wrapNodes, $isAtNodeEnd} from '@lexical/selection';
export { $getNearestNodeOfType, mergeRegister } from '@lexical/utils';
export { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext'

import LexicalErrorBoundary from '@lexical/react/LexicalErrorBoundary'

export const ErrorBoundary = LexicalErrorBoundary;