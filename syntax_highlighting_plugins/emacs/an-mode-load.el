;;; an-mode-load.el --- automatically extracted autoloads
;;; Commentary:

;; To install an-mode, add the following lines to your .emacs file:
;;   (add-to-list 'load-path "PATH CONTAINING an-mode-load.el" t)
;;   (require 'an-mode-load)
;;
;; After this, an-mode will be used for files ending in '.an'.
;;
;; To compile an-mode from the command line, run the following
;;   emacs -batch -f batch-byte-compile an-mode.el
;;
;; See an-mode.el for documentation.
;;
;; To update this file, evaluate the following form
;;   (let ((generated-autoload-file buffer-file-name)) (update-file-autoloads "an-mode.el"))

;;; Code:


;;;### (autoloads (an-download-play godoc gofmt-before-save an-mode)
;;;;;;  "an-mode" "an-mode.el" (20767 50749))
;;; Generated autoloads from an-mode.el

(autoload 'an-mode "an-mode" "\
Major mode for editing Antha source text.

This mode provides (not just) basic editing capabilities for
working with Antha code. It offers almost complete syntax
highlighting, indentation that is almost identical to anthafmt,
proper parsing of the buffer content to allow features such as
navigation by function, manipulation of comments or detection of
strings.

\(fn)" t nil)

(add-to-list 'auto-mode-alist (cons "\\.an\\'" 'an-mode))


(provide 'an-mode-load)
;; Local Variables:
;; version-control: never
;; no-byte-compile: t
;; no-update-autoloads: t
;; coding: utf-8
;; End:
;;; an-mode-load.el ends here
