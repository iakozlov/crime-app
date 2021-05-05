import pandas as pd
import numpy as np
import datetime
import calendar
import pickle
from joblib import dump, load
from sklearn.ensemble import RandomForestClassifier
import sys

def parse_date(s):
    res = {}
    dates = s.split()
    year = int((dates[0].split('-'))[0])
    res['year'] = year
    #years.append(year)
    month = int((dates[0].split('-'))[1])
    #months.append(month)
    res['month'] = month
    day = int((dates[0].split('-'))[2])
    res['day'] = day
    #days.append(day)
    hour = int((dates[1].split(':'))[0])
    res['hour'] = hour
    #hours.append(hour)
    minute = int((dates[1].split(':'))[1])
    res['minutes'] = minute
    #minutes.append(minute)
    if(minute == 30 or minute == 0):
        #beautiful_endings.append(1)
        res['beautiful_endings'] = 1
    else:
        res['beautiful_endings'] = 0
        #beautiful_endings.append(0)
    weekday = calendar.day_abbr[datetime.date(year, month, day).weekday()]
    res['weekday'] = weekday
    if(weekday == 'Sun' or weekday == 'Sat'):
        res['is_weekends'] = 1
        #is_weekends.append(1)
    else:
        res['is_weekends'] = 0
        #is_weekends.append(0)
    if(hour > 22 or (hour < 6)):
        res['nights'] = 1
        res['mornings'] = 0
        res['middle_days'] = 0
        res['afternoons'] = 0
    elif(hour >=6 and hour <= 10):
        res['nights'] = 0
        res['mornings'] = 1
        res['middle_days'] = 0
        res['afternoons'] = 0
    elif(hour > 10 and hour <=17):
        res['nights'] = 0
        res['mornings'] = 0
        res['middle_days'] = 1
        res['afternoons'] = 0
    else:
        res['nights'] = 0
        res['mornings'] = 0
        res['middle_days'] = 0
        res['afternoons'] = 1

    return res




def make_zones(X, Y, maxX = -122.365240723693, maxY = 37.819975492297004,
               minX = -122.513642064265, minY = 37.7078790224135):
    x_spread = maxX - minX
    y_spread = maxY - minY
    x_step = x_spread/15
    y_step = y_spread/15
    x_zone = int((X - minX)/x_step)
    y_zone = int((Y - minY)/y_step)
    if(x_zone == 15):
        x_zone = 14
    if(y_zone == 15):
        y_zone = 14
    return (15*x_zone+y_zone)

def parse_info(X, Y, date, district_clf, clf, address = ""):
    columns = np.array(['X','Y','BeautifulEndings','IsWeekend','Nights','Mornings','MiddleDays','Afternoons',
                        'Day','Months_1','Months_2','Months_3','Months_4','Months_5','Months_6','Months_7',
                        'Months_8','Months_9','Months_10','Months_11','Months_12','Hours','DayOfWeek_Friday',
                        'DayOfWeek_Monday','DayOfWeek_Saturday','DayOfWeek_Sunday','DayOfWeek_Thursday',
                        'DayOfWeek_Tuesday', 'DayOfWeek_Wednesday', 'PdDistrict_BAYVIEW', 'PdDistrict_CENTRAL',
                        'PdDistrict_INGLESIDE', 'PdDistrict_MISSION', 'PdDistrict_NORTHERN', 'PdDistrict_PARK',
                         'PdDistrict_RICHMOND', 'PdDistrict_SOUTHERN', 'PdDistrict_TARAVAL', 'PdDistrict_TENDERLOIN',
                         'Zones_2', 'Zones_3', 'Zones_4', 'Zones_5', 'Zones_6', 'Zones_7', 'Zones_8', 'Zones_9',
                         'Zones_15', 'Zones_16', 'Zones_17', 'Zones_18', 'Zones_19', 'Zones_20', 'Zones_21',
                         'Zones_22', 'Zones_23', 'Zones_24', 'Zones_25', 'Zones_30', 'Zones_31', 'Zones_32',
                         'Zones_33', 'Zones_34', 'Zones_35', 'Zones_36', 'Zones_37', 'Zones_38', 'Zones_39',
                         'Zones_40', 'Zones_41', 'Zones_45', 'Zones_46', 'Zones_47', 'Zones_48', 'Zones_49',
                         'Zones_50', 'Zones_51', 'Zones_52', 'Zones_53', 'Zones_54', 'Zones_55', 'Zones_58',
                         'Zones_60', 'Zones_61', 'Zones_62', 'Zones_63', 'Zones_64', 'Zones_65', 'Zones_66',
                         'Zones_67', 'Zones_68', 'Zones_69', 'Zones_70', 'Zones_75', 'Zones_76', 'Zones_77',
                         'Zones_78', 'Zones_79', 'Zones_80', 'Zones_81', 'Zones_82', 'Zones_83', 'Zones_84',
                         'Zones_85', 'Zones_86', 'Zones_87', 'Zones_90', 'Zones_91', 'Zones_92', 'Zones_93',
                         'Zones_94', 'Zones_95', 'Zones_96', 'Zones_97', 'Zones_98', 'Zones_99', 'Zones_100',
                         'Zones_101', 'Zones_102', 'Zones_103', 'Zones_105', 'Zones_106', 'Zones_107', 'Zones_108',
                         'Zones_109', 'Zones_110', 'Zones_111', 'Zones_112', 'Zones_113', 'Zones_114', 'Zones_115',
                         'Zones_116', 'Zones_117', 'Zones_118', 'Zones_120', 'Zones_121', 'Zones_122', 'Zones_123',
                         'Zones_124', 'Zones_125', 'Zones_126', 'Zones_127', 'Zones_128', 'Zones_129', 'Zones_130',
                         'Zones_131', 'Zones_132', 'Zones_133', 'Zones_135', 'Zones_136', 'Zones_137', 'Zones_138',
                         'Zones_139', 'Zones_140', 'Zones_141', 'Zones_142', 'Zones_143', 'Zones_144', 'Zones_145',
                         'Zones_146', 'Zones_147', 'Zones_148', 'Zones_150', 'Zones_151', 'Zones_152', 'Zones_153',
                         'Zones_154', 'Zones_155', 'Zones_156', 'Zones_157', 'Zones_158', 'Zones_159', 'Zones_160',
                         'Zones_161', 'Zones_162', 'Zones_163', 'Zones_165', 'Zones_166', 'Zones_167', 'Zones_168',
                         'Zones_169', 'Zones_170', 'Zones_171', 'Zones_172', 'Zones_173', 'Zones_174', 'Zones_175',
                         'Zones_176', 'Zones_177', 'Zones_178', 'Zones_180', 'Zones_181', 'Zones_182', 'Zones_183',
                         'Zones_184', 'Zones_185', 'Zones_186', 'Zones_187', 'Zones_188', 'Zones_189', 'Zones_190',
                         'Zones_191', 'Zones_195', 'Zones_196', 'Zones_197', 'Zones_198', 'Zones_199', 'Zones_200',
                         'Zones_201', 'Zones_212', 'Zones_213', 'Zones_223', 'Zones_224', 'XY1', 'XY2', 'XY3',
                         'XY4', 'XY45_2', 'XY30_1', 'XY30_2', 'XY60_1', 'XY60_2', 'XY5', 'XY_rad' ])

    arr = np.zeros(224)
    line = pd.DataFrame([arr], columns = columns)
    d = {}
    res = parse_date(date)
    line['X'][0] = X
    line['Y'][0] = Y

    line['BeautifulEndings'][0] = res['beautiful_endings']

    line['IsWeekend'][0] = res['is_weekends']

    line['Nights'][0] = res['nights']

    line['Mornings'][0] = res['mornings']

    line['MiddleDays'][0] = res['middle_days']

    line['Afternoons'][0] = res['afternoons']

    line['Day'][0] = res['day']

    line['Months_' + str(res['month'])] = 1

    line['Hours'] = res['hour']


    if(res['weekday'] == 'Fri'):
        line['DayOfWeek_Friday'][0] = 1
    elif(res['weekday'] == 'Mon'):
        line['DayOfWeek_Monday'][0] = 1
    elif(res['weekday'] == 'Sat'):
        line['DayOfWeek_Saturday'][0] = 1
    elif(res['weekday'] == 'Sun'):
        line['DayOfWeek_Sunday'][0] = 1
    elif(res['weekday'] == 'Thu'):
        line['DayOfWeek_Thursday'][0] = 1
    elif(res['weekday'] == 'Tue'):
        line['DayOfWeek_Tuesday'][0] = 1
    elif(res['weekday'] == 'Wed'):
        line['DayOfWeek_Wednesday'][0] = 1

    district = district_clf.predict([[X, Y]])
    if(district == 'BAYVIEW'):
        #arr[29] = 1
        line['PdDistrict_BAYVIEW'][0] = 1
    elif(district == 'CENTRAL'):
        #arr[30] = 1
        line['PdDistrict_CENTRAL'][0] = 1
    elif(district == 'INGLESIDE'):
        #arr[31] = 1
        line['PdDistrict_INGLESIDE'][0] = 1
    elif(district == 'MISSION'):
        #arr[32] = 1
        line['PdDistrict_MISSION'][0] = 1
    elif(district == 'PARK'):
        #arr[33] = 1
        line['PdDistrict_PARK'][0] = 1
    elif(district == 'RICHMOND'):
        #arr[34] = 1
        line['PdDistrict_RICHMOND'][0] = 1
    elif(district == 'SOUTHERN'):
        #arr[35] = 1
        line['PdDistrict_SOUTHERN'][0] = 1
    elif(district == 'TARAVAL'):
        #arr[36] = 1
        line['PdDistrict_TARAVAL'][0] = 1
    elif(district == 'TENDERLOIN'):
        #arr[37] = 1
        line['PdDistrict_TENDERLOIN'][0] = 1


    zone = make_zones(X, Y)
    if('Zones_' + str(zone) in columns):
        line['Zones_' + str(zone)] = 1

    maxX = -122.365240723693
    maxY = 37.819975492297004
    minX = -122.513642064265
    minY = 37.7078790224135

    line['XY1'][0] = (X - minX)**2 + (Y - minY)**2
    line['XY2'][0] = (maxX - X)**2 + (Y - minY)**2
    line['XY3'][0] = (X - minX)**2 + (maxY - Y)**2
    line['XY4'][0] = (maxX - X)**2 + (maxY - Y)**2
    line['XY45_2'][0] = Y * np.cos(np.pi / 4) - X * np.sin(np.pi / 4)
    line['XY30_1'][0] = X * np.cos(np.pi / 6) + Y * np.sin(np.pi / 6)
    line['XY30_2'][0] = Y * np.cos(np.pi / 6) - X * np.sin(np.pi / 6)
    line['XY60_1'] = X * np.cos(np.pi / 3) + Y * np.sin(np.pi / 3)

    line['XY60_2'] = Y * np.cos(np.pi / 3) - X * np.sin(np.pi / 3)

    X_median = -122.416452065595
    Y_median = 37.775420706711

    line['XY5'][0] = (X - X_median) ** 2 + (Y - Y_median) ** 2

    line['XY_rad'][0] = np.sqrt(np.power(Y, 2) + np.power(X, 2))
    d = {}
    predicted = clf.predict_proba(line)[0]
    for i in range(len(clf.classes_)):
        d[clf.classes_[i]] = predicted[i]
    sorted_dict = {}
    sorted_keys = sorted(d, key=d.get)

    for w in sorted_keys:
        sorted_dict[w] = d[w]
    s = ''
    keys = list(sorted_dict.keys())
    #display(keys)
    values = list(sorted_dict.values())
    for i in range(3):
        s += str(keys[len(keys) - 1 - i]) + ':' + str(values[len(keys) - 1 - i]) + ';'
    #display(sorted_dict)

    return s


clf = load('model.joblib')
district_clf = load('district_model.joblib')
arguments = sys.argv[1:]
f = open("info.txt", "w")
f.write(parse_info(float(arguments[0]), float(arguments[1]), arguments[2], district_clf, clf))
f.close()
