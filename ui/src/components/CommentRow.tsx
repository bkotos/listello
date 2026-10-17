import type { ItemCommentDto } from "api-types";
import { formatDistanceToNow } from "date-fns";
import { Trash2 } from "lucide-react";

export type CommentRowProps = {
  comment: ItemCommentDto;
  currentUserId?: string;
  onDelete?: (commentId: string) => void;
};

export function CommentRow({ comment, currentUserId = "", onDelete }: CommentRowProps) {
  return (
    <div className="hover-parent is-flex" style={{ gap: "0.625rem" }}>
      <div
        className="is-flex is-align-items-center is-justify-content-center is-family-code has-text-primary has-text-weight-semibold"
        style={{
          height: "1.75rem",
          width: "1.75rem",
          flexShrink: 0,
          borderRadius: "9999px",
          fontSize: "0.6875rem",
          backgroundColor:
            "hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l) / 0.15)",
        }}
      >
        {comment.UserName.slice(0, 1)}
      </div>
      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <div className="is-flex is-align-items-center" style={{ gap: "0.5rem" }}>
          <span className="is-size-7 has-text-weight-medium">{comment.UserName}</span>
          <span className="is-size-7 muted">{relativeTime(comment.CreatedAt)}</span>
          {comment.UserID === currentUserId ? (
            <button
              type="button"
              aria-label="Delete comment"
              onClick={() => onDelete?.(comment.ID)}
              className="icon-btn is-danger-hover hover-reveal"
              style={{ marginLeft: "auto", height: "1.5rem", width: "1.5rem" }}
            >
              <Trash2 size={14} />
            </button>
          ) : null}
        </div>
        <p className="is-size-7 muted">{comment.Body}</p>
      </div>
    </div>
  );
}

function relativeTime(createdAt: string): string {
  const date = new Date(createdAt);
  if (Number.isNaN(date.getTime())) {
    return "";
  }
  return formatDistanceToNow(date, { addSuffix: true });
}
