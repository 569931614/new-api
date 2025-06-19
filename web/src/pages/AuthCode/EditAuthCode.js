import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Card,
  DatePicker,
  Form,
  Input,
  Modal,
  Select,
  Switch,
  TextArea,
  Typography,
  Space,
  Divider,
  Tag
} from '@douyinfe/semi-ui';
import { IconShield, IconUser, IconClock, IconCode, IconSetting } from '@douyinfe/semi-icons';
import { API, showError, showSuccess } from '../../helpers';

const { Title } = Typography;

const EditAuthCode = (props) => {
  const { t } = useTranslation();
  const originInputs = {
    code: '',
    name: '',
    description: '',
    user_type: 1,
    expired_time: -1,
    is_bot: false,
    wx_auto_x_code: '',
    machine_code: '',
    group: '',
    status: 1
  };

  const [inputs, setInputs] = useState(originInputs);
  const [loading, setLoading] = useState(false);
  const [expiredDate, setExpiredDate] = useState(null);
  const [formApi, setFormApi] = useState(null);
  const [dataLoaded, setDataLoaded] = useState(false);
  const [groupOptions, setGroupOptions] = useState([]);

  const handleInputChange = (name, value) => {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  };

  const loadAuthCode = async () => {
    try {
      let res = await API.get(`/api/auth_code/${props.editingAuthCode.id}`);
      const { success, message, data } = res.data;
      if (success) {
        setInputs(data);
        if (data.expired_time !== -1) {
          setExpiredDate(new Date(data.expired_time * 1000));
        } else {
          setExpiredDate(null);
        }
        setDataLoaded(true);
      } else {
        showError(message);
      }
    } catch (error) {
      showError(error.message);
    }
  };

  const fetchGroups = async () => {
    try {
      let res = await API.get(`/api/group/`);
      if (res === undefined) {
        return;
      }
      setGroupOptions(
        res.data.data.map((group) => ({
          label: group,
          value: group,
        })),
      );
    } catch (error) {
      showError(error.message);
    }
  };

  useEffect(() => {
    // 只在弹窗打开时处理数据
    if (props.visible) {
      // 加载分组数据
      fetchGroups();

      if (props.editingAuthCode.id) {
        setDataLoaded(false);
        loadAuthCode().then();
      } else {
        setInputs(originInputs);
        setExpiredDate(null);
        setDataLoaded(true);
      }
    }
  }, [props.visible, props.editingAuthCode.id]);



  const submit = async () => {
    if (!inputs.code || !inputs.name) {
      showError(t('请填写授权码和名称！'));
      return;
    }

    setLoading(true);

    // 处理过期时间
    let submitInputs = { ...inputs };
    if (expiredDate) {
      submitInputs.expired_time = Math.floor(expiredDate.getTime() / 1000);
    } else {
      submitInputs.expired_time = -1;
    }

    try {
      let res;
      if (props.editingAuthCode.id) {
        res = await API.put('/api/auth_code/', submitInputs);
      } else {
        res = await API.post('/api/auth_code/', submitInputs);
      }
      const { success, message } = res.data;
      if (success) {
        if (props.editingAuthCode.id) {
          showSuccess(t('授权码更新成功！'));
        } else {
          showSuccess(t('授权码创建成功！'));
          setInputs(originInputs);
          setExpiredDate(null);
        }
        props.refresh();
        props.onCancel();
      } else {
        showError(message);
      }
    } catch (error) {
      showError(error.message);
    }
    setLoading(false);
  };

  const userTypeOptions = [
    { label: t('普通用户'), value: 1 },
    { label: t('管理员'), value: 10 },
    { label: t('超级管理员'), value: 100 }
  ];

  const statusOptions = [
    { label: t('启用'), value: 1 },
    { label: t('禁用'), value: 2 },
    { label: t('已使用'), value: 3 },
    { label: t('待激活'), value: 4 },
    { label: t('激活'), value: 5 }
  ];

  return (
    <Modal
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <IconShield style={{ color: '#059669', fontSize: '20px' }} />
          <Title heading={4} style={{ margin: 0 }}>
            {props.editingAuthCode.id ? t('编辑授权码') : t('创建授权码')}
          </Title>
        </div>
      }
      visible={props.visible}
      onCancel={props.onCancel}
      onOk={submit}
      okText={t('提交')}
      cancelText={t('取消')}
      confirmLoading={loading}
      width={700}
      centered={true}
      bodyStyle={{ padding: 0, maxHeight: '80vh', overflow: 'auto' }}
      style={{ borderRadius: '12px' }}
    >
      <div style={{ padding: '20px' }}>
        <div style={{ background: 'transparent' }}>
          {/* 基本信息区域 */}
          <div style={{ marginBottom: '20px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '12px',
              color: '#1890ff',
              fontWeight: 600,
              fontSize: '15px'
            }}>
              <IconCode />
              <span>{t('基本信息')}</span>
            </div>

            {/* 使用两列布局 */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '12px' }}>
              <div>
                <div style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  <span>{t('授权码')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
                <Input
                  placeholder={t('请输入授权码')}
                  value={inputs.code || ''}
                  onChange={(value) => setInputs(prev => ({ ...prev, code: value }))}
                  style={{ width: '100%' }}
                />
              </div>

              <div>
                <div style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  <span>{t('名称')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
                <Input
                  placeholder={t('请输入名称')}
                  value={inputs.name || ''}
                  onChange={(value) => setInputs(prev => ({ ...prev, name: value }))}
                  style={{ width: '100%' }}
                />
              </div>
            </div>

            <div style={{ marginBottom: '12px' }}>
              <div style={{
                marginBottom: '6px',
                fontSize: '13px',
                fontWeight: 500
              }}>
                {t('描述')}
              </div>
              <TextArea
                placeholder={t('请输入描述信息')}
                value={inputs.description || ''}
                onChange={(value) => setInputs(prev => ({ ...prev, description: value }))}
                autosize={{ minRows: 2, maxRows: 3 }}
                style={{ width: '100%' }}
              />
            </div>
          </div>

          <Divider margin="12px" />

          {/* 权限配置区域 */}
          <div style={{ marginBottom: '20px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '12px',
              color: '#52c41a',
              fontWeight: 600,
              fontSize: '15px'
            }}>
              <IconUser />
              <span>{t('权限配置')}</span>
            </div>

            {/* 使用两列布局 */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', alignItems: 'end' }}>
              <div>
                <div style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  <span>{t('用户类型')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
                <Select
                  value={inputs.user_type}
                  onChange={(value) => setInputs(prev => ({ ...prev, user_type: value }))}
                  optionList={userTypeOptions}
                  style={{ width: '100%' }}
                />
              </div>

              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <Switch
                  checked={inputs.is_bot}
                  onChange={(checked) => setInputs(prev => ({ ...prev, is_bot: checked }))}
                />
                <span style={{ fontSize: '13px', fontWeight: 500 }}>{t('机器人账户')}</span>
              </div>
            </div>
          </div>

          <Divider margin="12px" />

          {/* 时间配置和高级配置合并 */}
          <div style={{ marginBottom: '20px' }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              marginBottom: '12px',
              color: '#fa8c16',
              fontWeight: 600,
              fontSize: '15px'
            }}>
              <IconClock />
              <span>{t('时间与高级配置')}</span>
            </div>

            {/* 使用两列布局 */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginBottom: '12px' }}>
              <div>
                <div style={{
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  {t('过期时间')}
                </div>
                <DatePicker
                  placeholder={t('选择过期时间')}
                  value={expiredDate}
                  onChange={(date) => setExpiredDate(date)}
                  type="dateTime"
                  style={{ width: '100%' }}
                />
              </div>

              <div>
                <div style={{
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  {t('WxAutoX码')}
                </div>
                <Input
                  placeholder={t('请输入WxAutoX码（可选）')}
                  value={inputs.wx_auto_x_code || ''}
                  onChange={(value) => setInputs(prev => ({ ...prev, wx_auto_x_code: value }))}
                  style={{ width: '100%' }}
                />
              </div>
            </div>

            {/* 机器码字段 */}
            <div style={{ marginBottom: '12px' }}>
              <div style={{
                marginBottom: '6px',
                fontSize: '13px',
                fontWeight: 500
              }}>
                {t('机器码')}
              </div>
              <Input
                placeholder={t('请输入机器码（可选）')}
                value={inputs.machine_code || ''}
                onChange={(value) => setInputs(prev => ({ ...prev, machine_code: value }))}
                style={{ width: '100%' }}
                disabled={props.editingAuthCode.id && inputs.machine_code && inputs.status === 5}
              />
              {props.editingAuthCode.id && inputs.machine_code && inputs.status === 5 && (
                <div style={{
                  fontSize: '11px',
                  color: '#52c41a',
                  marginTop: '4px'
                }}>
                  {t('机器码已绑定并激活，无法修改')}
                </div>
              )}
            </div>

            {/* 分组字段 */}
            <div style={{ marginBottom: '12px' }}>
              <div style={{
                marginBottom: '6px',
                fontSize: '13px',
                fontWeight: 500
              }}>
                {t('分组')}
              </div>
              <Select
                placeholder={t('请选择分组（可选）')}
                value={inputs.group ? inputs.group.split(',').filter(g => g.trim() !== '') : []}
                onChange={(value) => {
                  const groupStr = Array.isArray(value) ? value.join(',') : '';
                  setInputs(prev => ({ ...prev, group: groupStr }));
                }}
                optionList={groupOptions}
                multiple
                allowAdditions
                additionLabel={t('请在系统设置页面编辑分组倍率以添加新的分组：')}
                style={{ width: '100%' }}
                showClear
              />
              <div style={{
                fontSize: '11px',
                color: '#666',
                marginTop: '4px'
              }}>
                {t('设置分组后，可通过外部接口获取该分组下的渠道列表')}
              </div>
            </div>

            {props.editingAuthCode.id && (
              <div style={{ maxWidth: '50%' }}>
                <div style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '4px',
                  marginBottom: '6px',
                  fontSize: '13px',
                  fontWeight: 500
                }}>
                  <span>{t('状态')}</span>
                  <Tag color="red" size="small">{t('必填')}</Tag>
                </div>
                <Select
                  value={inputs.status}
                  onChange={(value) => setInputs(prev => ({ ...prev, status: value }))}
                  optionList={statusOptions}
                  style={{ width: '100%' }}
                />
              </div>
            )}

            <div style={{
              fontSize: '11px',
              color: '#8c8c8c',
              marginTop: '8px'
            }}>
              {t('不设置过期时间表示永久有效')}
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default EditAuthCode;
