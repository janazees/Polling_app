import "./GradientField.css";

/**
 * Decorative background: a slow-moving gradient wash (butter -> sage -> umber)
 * under a graph-paper grid, matching the Figma reference. Purely presentational
 * — aria-hidden so screen readers skip it.
 */
export default function GradientField() {
  return (
    <div className="gradient-field" aria-hidden="true">
      <div className="gradient-field__wash" />
      <div className="gradient-field__grid" />
    </div>
  );
}
