define([
	'jquery',
	'base',
	'dialog'
], function($, lotto) {
	"use strict";

	lotto.EditDialog = function(elementId, parent) {
		var that = this,
			lang = this._lang;
		lotto.Dialog.prototype.constructor.call(this, elementId, parent);
		this.lottoId = elementId + '_lotto';
		this._lottoId = '#' + elementId + '_lotto'
		this.plusId = elementId + '_plus';
		this._plusId = '#' + elementId + '_plus';
		this.startDateId = elementId + '_start_date';
		this._startDateId = '#' + elementId + '_start_date';
		this.endDateId = elementId + '_end_date';
		this._endDateId = '#' + elementId + '_end_date';
		this.buttonId = elementId + '_btn';
		this._buttonId = '#' + elementId + '_btn';
		this.buttonCloseId = ''+elementId+'close_btn';
		this._buttonCloseId = '#'+elementId+'close_btn';
		
		this._options = {
			modal: true,
			width: 600,
			title: this._lang.bt_add,
			position: {
          		my: "center top",	
				at: "center top",	
        	},
			open: function() {
      			$(":button:contains('Close')").focus(); // Set focus to the [Ok] button
    		},
			buttons: [
				{
					text: "Save",
					click: that.okClick.bind(that)
				},
				{
					text: "Close",
					click: that.close.bind(that)
				}
			]
		};
		this.lottoRE = /^\d{1,2}\s\d{1,2}\s\d{1,2}\s\d{1,2}\s\d{1,2}\s\d{1,2}\s*$/;
		this.dateRE = [/^(\d{1,2})-(\d{1,2})-(\d{4})$/,  
						/^(\d{1,2})-(\d{1,2})-(\d{2})$/,
						/^(\d{1,2})-(\d{1,2})$/]
	}
	$.extend(lotto.EditDialog.prototype, lotto.Dialog.prototype, {
		_getHTML: function() {
			var lang = this._lang;
			return ['<div id="', this.elementId, '">',
				'<table><tr>',
				'<tr><td>Lotto</td><td><input id="', this.lottoId, '"/></td></tr>',
				'<tr><td>Plus</td><td><input type="checkbox" id="', this.plusId, '"/></td></tr>',
				'<tr><td>', lang.startDate, '</td><td><input id="', this.startDateId, '"/></td></tr>',
				'<tr><td>', lang.endDate, '</td><td><input id="', this.endDateId, '"/></td></tr>',
				'</table>',
				'</div > '
			].join('');
		},
		_postInit: function() {
			$(this._buttonId).button().on('click', this.okClick.bind(this));
			$(this._buttonCloseId).button().on('click', this.close.bind(this));
			
		},
		okClick: function() {
			var data = {
				lotto: $(this._lottoId).val(),
				plus: $(this._plusId).is(':checked'),
				startDate: $(this._startDateId).val(),
				endDate: $(this._endDateId).val()
			};
			if (!this.validate(data)) {
				return;
			}
			$.post('/saveLucky', data, function(data) {
				this.saved(data)
			}.bind(this));
			
		},
		validate: function(data) {
			if (!this.lottoRE.exec(data.lotto)) {
				alert(this._lang.wrongLottoFortmat);
				return false;
			}
			if (!this.validateDate(data.startDate)) {
				alert(this._lang.wrongStartDateFortmat);
				return false;
			}
			if (!this.validateDate(data.endDate)) {
				alert(this._lang.wrongEndDateFortmat);
				return false;
			}
			return true;
		},
		validateDate: function(value) {
			for (var i=0; i < this.dateRE.length; i++){
				var date = this.dateRE[i].exec(value);
				if (!date) {
					continue;
				}
				var month = +date[2];
				if (month < 1 || month > 12) {
					return false;
				}
				var day = +date[1];
				if (day > 31 || day < 1) {
					return false;
				}
				return true;
			}
		},
		saved: function(data) {
			var obj = JSON.parse(data);
			if (obj.status != "ok") {
				alert(obj.msg);
			}else{
				$.getJSON('/getLucky').done(this.update.bind(this));
			}
		},
		open: function(){
			var that = this;
			$.getJSON('/getLucky', function(data){
				$(that._lottoId).val(data.lucky);
				$(that._plusId).prop("checked", data.plus);
				$(that._startDateId).val(data.startDate);
				$(that._endDateId).val(data.endDate);
			});
			lotto.Dialog.prototype.open.call(this);
		},
		update: function(data){
			this._app.updateLucky(data);
			this.close();
		}
	});
	return lotto.EditDialog;
});