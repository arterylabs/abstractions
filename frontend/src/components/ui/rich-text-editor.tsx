"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import { TextSelection } from "@tiptap/pm/state";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import Tooltip from "@/components/ui/tooltip";
import { EmojiExtension } from "@/lib/emoji-extension";
import { CustomImage } from "@/lib/image-extension";
import { Notifier } from "@/components/ui/sonner";
import Icon from "@/components/ui/icon";
import filesApi from "@/api/files";
import { getErrorMessage } from "@/lib/errors";

interface RichTextEditorProps {
  value?: string;
  onChange?: (value: string) => void;
  placeholder?: string;
  className?: string;
  editorClassName?: string;
  /** When false the editor renders as read-only and the toolbar is hidden. Defaults to true. */
  editable?: boolean;
  /**
   * When true, disables all formatting — no toolbar is shown, paste is stripped
   * to plain text, and onChange emits plain text instead of HTML.
   */
  formattingDisabled?: boolean;
}

const ToolbarButton = ({
  onClick,
  isActive,
  children,
  disabled,
  tooltip,
}: {
  onClick: () => void;
  isActive?: boolean;
  children: React.ReactNode;
  disabled?: boolean;
  tooltip: string;
}) => (
  <Tooltip value={tooltip}>
    <Button
      type="button"
      variant="ghost"
      size="sm"
      className={cn(
        "h-7 w-7 p-0 font-medium",
        isActive &&
          "bg-primary/15 text-primary hover:bg-primary/20 hover:text-primary",
      )}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </Button>
  </Tooltip>
);

export function RichTextEditor({
  value,
  onChange,
  placeholder = "Start writing…",
  className,
  editorClassName,
  editable = true,
  formattingDisabled = false,
}: RichTextEditorProps) {
  const [uploadingImageCount, setUploadingImageCount] = useState(0);
  const fileInputRef = useRef<HTMLInputElement>(null);
  // Stable ref to the tiptap editor used inside editorProps.handlePaste (which is
  // created before the editor instance is returned from useEditor).
  const editorRef = useRef<ReturnType<typeof useEditor>>(null);

  /**
   * Uploads a File object as a public image and inserts it at the current
   * cursor position. Returns true on success, false on failure.
   */
  const replaceUploadingImage = useCallback(
    (
      editorInstance: NonNullable<ReturnType<typeof useEditor>>,
      uploadId: string,
      src?: string,
    ) => {
      let targetPos: number | null = null;
      let targetNodeSize = 0;
      let targetAttrs: Record<string, unknown> | null = null;

      editorInstance.state.doc.descendants((node, pos) => {
        if (node.type.name === "image" && node.attrs.uploadId === uploadId) {
          targetPos = pos;
          targetNodeSize = node.nodeSize;
          targetAttrs = node.attrs as Record<string, unknown>;
          return false;
        }
        return true;
      });

      if (targetPos === null) return;
      const tr = editorInstance.state.tr;

      if (src) {
        tr.setNodeMarkup(targetPos, undefined, {
          ...(targetAttrs || {}),
          src,
          uploadId: null,
          isUploading: null,
          alt: null,
          title: null,
        });
      } else {
        tr.delete(targetPos, targetPos + targetNodeSize);
      }

      editorInstance.view.dispatch(tr);
    },
    [],
  );

  const uploadAndInsertImage = useCallback(
    async (
      file: File,
      editorInstance: NonNullable<ReturnType<typeof useEditor>>,
      uploadId?: string,
    ) => {
      setUploadingImageCount((count) => count + 1);
      try {
        const formData = new FormData();
        formData.append("file", file);
        formData.append("public", "true");
        const res = await filesApi.uploadFile(formData);
        const link = res.data?.link;
        if (!link) throw new Error("No link returned from upload");
        if (uploadId) {
          replaceUploadingImage(editorInstance, uploadId, link);
        } else {
          editorInstance.chain().focus().setImage({ src: link }).run();
        }
        return true;
      } catch (error) {
        if (uploadId) {
          replaceUploadingImage(editorInstance, uploadId);
        }
        Notifier.error({
          title: "Failed to upload image",
          message: getErrorMessage(error),
        });
        return false;
      } finally {
        setUploadingImageCount((count) => Math.max(0, count - 1));
      }
    },
    [replaceUploadingImage],
  );

  const editor = useEditor({
    immediatelyRender: false,
    editable,
    extensions: formattingDisabled
      ? [
          // Plain-text mode: only core document/paragraph/text nodes + undo history.
          // All formatting marks and block-level nodes are excluded.
          StarterKit.configure({
            bold: false,
            italic: false,
            strike: false,
            code: false,
            codeBlock: false,
            heading: false,
            bulletList: false,
            orderedList: false,
            listItem: false,
            blockquote: false,
            horizontalRule: false,
          }),
          Placeholder.configure({ placeholder }),
        ]
      : [
          StarterKit,
          Placeholder.configure({ placeholder }),
          // Image extension: draggable (node spec has draggable:true), base64 disabled.
          // Selection ring is applied via .ProseMirror-selectednode CSS below.
          CustomImage,
          // Emoji autocomplete triggered by ':' — inserts Unicode emoji characters.
          EmojiExtension,
        ],
    content: value || "",
    onUpdate: ({ editor }) => {
      onChange?.(editor.getHTML());
    },
    editorProps: {
      attributes: {
        class: cn(
          "min-h-[120px] w-full px-3 py-2 text-sm text-foreground focus:outline-none",
          // Plain-text mode: reset font weight so content never appears bold.
          formattingDisabled && "font-normal",
          !formattingDisabled && [
            "[&_p]:mb-2 [&_h1]:text-2xl [&_h1]:font-bold [&_h1]:mb-2",
            "[&_h2]:text-xl [&_h2]:font-semibold [&_h2]:mb-2",
            "[&_ul]:list-disc [&_ul]:pl-5 [&_ol]:list-decimal [&_ol]:pl-5",
            "[&_blockquote]:border-l-4 [&_blockquote]:border-primary/40 [&_blockquote]:pl-4 [&_blockquote]:italic [&_blockquote]:text-muted-foreground [&_blockquote]:my-3",
            "[&_img]:max-w-full [&_img]:rounded-md [&_img]:my-2 [&_img]:block [&_img]:cursor-grab",
            editable &&
              "[&_img.ProseMirror-selectednode]:ring-2 [&_img.ProseMirror-selectednode]:ring-primary [&_img.ProseMirror-selectednode]:ring-offset-2 [&_img.ProseMirror-selectednode]:outline-none",
            "[&_img[data-uploading='true']]:w-full [&_img[data-uploading='true']]:min-h-40 [&_img[data-uploading='true']]:max-h-72 [&_img[data-uploading='true']]:object-cover [&_img[data-uploading='true']]:opacity-70 [&_img[data-uploading='true']]:animate-pulse [&_img[data-uploading='true']]:cursor-default",
          ],
          "[&_.tiptap_p.is-editor-empty:first-child_::before]:text-muted-foreground",
          editorClassName ?? "",
        ),
      },
      handlePaste(_view, event) {
        // In plain-text mode strip all formatting and insert as plain text.
        if (formattingDisabled) {
          const text = event.clipboardData?.getData("text/plain") ?? "";
          if (text) {
            event.preventDefault();
            editorRef.current?.commands.insertContent(text);
            return true;
          }
          return false;
        }
        if (!editable) return false;
        const items = Array.from(event.clipboardData?.items ?? []);
        const imageItem = items.find(
          (item) => item.kind === "file" && item.type.startsWith("image/"),
        );
        if (!imageItem) return false;
        const file = imageItem.getAsFile();
        if (!file) return false;
        // Upload async; tiptap paste is blocked for this item.
        // editorRef.current is populated by the time the user can paste.
        if (editorRef.current) {
          const uploadId = `upload-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
          const placeholderSvg = encodeURIComponent(
            `<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="560" viewBox="0 0 1200 560">
              <rect width="1200" height="560" fill="#f4f4f5"/>
              <g transform="translate(600 240)">
                <circle cx="0" cy="0" r="26" fill="none" stroke="#d4d4d8" stroke-width="8"/>
                <path d="M0 -26 A26 26 0 0 1 22 -13" fill="none" stroke="#71717a" stroke-width="8" stroke-linecap="round">
                  <animateTransform attributeName="transform" attributeType="XML" type="rotate" from="0 0 0" to="360 0 0" dur="0.8s" repeatCount="indefinite"/>
                </path>
              </g>
              <text x="600" y="312" text-anchor="middle" font-family="Arial, sans-serif" font-size="24" fill="#71717a">Uploading image…</text>
            </svg>`,
          );
          editorRef.current
            .chain()
            .focus()
            .insertContent({
              type: "image",
              attrs: {
                src: `data:image/svg+xml;charset=utf-8,${placeholderSvg}`,
                alt: "Uploading image",
                uploadId,
                isUploading: "true",
              },
            })
            .run();
          uploadAndInsertImage(file, editorRef.current, uploadId);
        }
        return true;
      },
    },
  });

  // Keep the stable ref in sync so handlePaste can access the editor instance.
  useEffect(() => {
    editorRef.current = editor;
  }, [editor]);

  // Sync editable prop to tiptap whenever it changes.
  // Pass `false` as second arg to suppress the update event tiptap fires by default,
  // which would otherwise trigger onUpdate → onChange → setIsDirty spuriously.
  useEffect(() => {
    if (editor && editor.isEditable !== editable) {
      editor.setEditable(editable, false);
    }
  }, [editor, editable]);

  useEffect(() => {
    if (!editor || editable) return;
    const selection = TextSelection.atStart(editor.state.doc);
    const tr = editor.state.tr.setSelection(selection);
    editor.view.dispatch(tr);
    editor.view.dom.blur();
  }, [editor, editable]);

  // Sync external value changes (e.g. when the parent resets content on cancel).
  // emitUpdate: false prevents onUpdate from firing, so this
  // won't spuriously mark the form as dirty.
  useEffect(() => {
    if (!editor) return;
    const current = editor.getHTML();
    if (current !== value) {
      editor.commands.setContent(value || "", { emitUpdate: false });
    }
  }, [editor, value]);

  const handleImageFileChange = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file || !editor) return;
      // Reset the input so the same file can be re-selected later
      e.target.value = "";
      await uploadAndInsertImage(file, editor);
    },
    [editor, uploadAndInsertImage],
  );

  if (!editor) return null;

  return (
    <div
      className={cn(
        "rounded-md border border-input bg-background shadow-xs transition-[color,box-shadow]",
        "ring-ring/10 dark:ring-ring/20 outline-ring/50 dark:outline-ring/40",
        "focus-within:ring-4 focus-within:outline-1 !border-solid",
        // Read-only mode: strip border, shadow and focus ring so it renders as plain content.
        !editable &&
          "border-0 shadow-none ring-0 outline-none focus-within:ring-0 focus-within:outline-none bg-transparent",
        className,
      )}
    >
      {/* Toolbar — only shown in edit mode with formatting enabled */}
      {editable && !formattingDisabled && (
        <div className="flex flex-wrap gap-1 border-b border-input px-2 py-1">
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBold().run()}
            isActive={editor.isActive("bold")}
            tooltip="Bold"
          >
            <Icon type="format_bold" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleItalic().run()}
            isActive={editor.isActive("italic")}
            tooltip="Italic"
          >
            <Icon type="format_italic" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() =>
              editor.chain().focus().toggleHeading({ level: 2 }).run()
            }
            isActive={editor.isActive("heading", { level: 2 })}
            tooltip="Heading 2"
          >
            <Icon type="format_h2" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() =>
              editor.chain().focus().toggleHeading({ level: 3 }).run()
            }
            isActive={editor.isActive("heading", { level: 3 })}
            tooltip="Heading 3"
          >
            <Icon type="format_h3" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBulletList().run()}
            isActive={editor.isActive("bulletList")}
            tooltip="Bullet list"
          >
            <Icon type="format_list_bulleted" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
            isActive={editor.isActive("orderedList")}
            tooltip="Ordered list"
          >
            <Icon type="format_list_numbered" className="text-base" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBlockquote().run()}
            isActive={editor.isActive("blockquote")}
            tooltip="Blockquote"
          >
            <Icon type="format_quote" className="text-base" />
          </ToolbarButton>

          {/* Image upload button */}
          <ToolbarButton
            onClick={() => fileInputRef.current?.click()}
            disabled={uploadingImageCount > 0}
            tooltip="Insert image"
          >
            {uploadingImageCount > 0 ? (
              <Icon
                type="progress_activity"
                className="animate-spin text-base"
              />
            ) : (
              <Icon type="image" className="text-base" />
            )}
          </ToolbarButton>

          {/* Hidden file input — accepts all image formats */}
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleImageFileChange}
          />
        </div>
      )}
      <EditorContent editor={editor} />
    </div>
  );
}
