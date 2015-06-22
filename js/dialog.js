define([
	'jquery',
	'base'
], function($, lotto) {
	"use strict";
	lotto.Dialog = function(elementId, parent) {
		this._app = parent;
		this._lang = parent.lang;
		this._dialog = null;
		this.elementId = elementId;
		this._elementId = '#' + elementId;
		this._options = {};
	}
	$.extend(lotto.Dialog.prototype, {
		_getHTML: function() {
			return '';
		},
		_init: function() {
			if ($(this._elementId).length == 0) {
				$('body').append(this._getHTML());
				this._dialog = $(this._elementId);
				this._dialog.dialog(this._options);
				this._postInit();
			}
			return this._dialog;
		},
		_postInit: function(){

		},
		open: function(){
			if(this.isOpen()){
				return;
			}
			var dialog = this._init();
			dialog.dialog('open');
		},
		isOpen: function(){
			return (this._dialog != null) && this._dialog.dialog('isOpen');
		},
		close: function(){
			if(this.isOpen()){
				this._dialog.dialog('close');
			}
		}

	});
	return lotto.Dialog;
});