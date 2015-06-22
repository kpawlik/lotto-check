define([
	'jquery',
	'base',
	'lang/lang',
	'editDialog',
	'lang/en/lang',
	'lang/pl/lang'
], function($, lotto, lang) {
	"use strict";

	lotto.App = function(logoutUrl){
		this._logoutUrl = logoutUrl;
		this.lang = lang['en'];

	}
	$.extend(lotto.App.prototype, {
		init: function(){
			var lang = this.lang;
			$('#mi-edit').text(lang.bt_add).button().on('click', this.openEditDialog.bind(this));
			$('#mi-stats').text(lang.bt_history).button().on('click', this.gotoHistory.bind(this));
			$('#mi-logout').text(lang.bt_logout).button().on('click', this.logout.bind(this));
			$('#mi-main').text(lang.bt_main).button().on('click', this.gotoMain.bind(this));
			$('#user_info').text(lang.user_info);
		},

		openEditDialog: function(){
			if(!this.editDialog){
				this.editDialog = new lotto.EditDialog('lotto-edit-dialog', this);
			}
			this.editDialog.open();
		},
		logout: function(){
			window.location = this._logoutUrl;
		},
		updateLucky: function(data){
			$("#panel_start_date").html(data.startDate);
			$("#panel_end_date").html(data.endDate);
			$("#panel_lucky").html(data.lucky);
			$("#panel_plus").html(data.plus? 'Yes': 'No');
		},
		gotoHistory: function(){
			window.location = '/history'
		},
		gotoMain: function(){
			window.location = '/'
		}

	});
	return lotto;
});