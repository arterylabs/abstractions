import Image from "@tiptap/extension-image";

/**
 * Extends the official @tiptap/extension-image.
 * - `draggable: true` is already set in the base extension's node spec.
 * - Selection styling is handled via the `.ProseMirror-selectednode` CSS class
 *   that ProseMirror adds automatically to the selected node's DOM element —
 *   no ReactNodeViewRenderer needed, which avoids a DOM ownership conflict
 *   between ProseMirror and React when setContent() is called.
 */
export const CustomImage = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      uploadId: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-upload-id"),
        renderHTML: (attributes) => {
          if (!attributes.uploadId) return {};
          return { "data-upload-id": attributes.uploadId };
        },
      },
      isUploading: {
        default: null,
        parseHTML: (element) => element.getAttribute("data-uploading"),
        renderHTML: (attributes) => {
          if (!attributes.isUploading) return {};
          return { "data-uploading": attributes.isUploading };
        },
      },
    };
  },
}).configure({
  inline: false,
  // Disallow base64 — images must always be uploaded to the backend
  allowBase64: false,
});
