// 为周杰伦创建5个病历的浏览器控制台脚本
// 使用方法：在 http://localhost:8080 登录医生账号后，打开浏览器控制台，复制粘贴此脚本并回车

async function createPrescriptionsForZhouJielun() {
    const baseURL = 'http://localhost:8888';
    
    // 获取当前登录的token（从Cookie中获取）
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
    }
    
    const token = getCookie('account_id_token');
    if (!token) {
        console.error('❌ 未找到登录token，请先登录医生账号');
        console.log('提示：请确保已经登录到系统');
        return;
    }
    
    console.log('✅ 已获取token:', token)

    // 5个病历数据 - 3个温江社区，2个协和医院
    const prescriptions = [
        {
            hospital: '温江社区医疗中心',
            patientName: '周杰伦',
            diagnosis: '急性上呼吸道感染',
            symptoms: '咽痛、流涕、轻度发热（37.8℃）、咳嗽',
            treatment: '对症治疗，多休息，多饮水',
            medications: '阿莫西林胶囊 0.5g 每日3次；布洛芬缓释胶囊 0.3g 每日2次；复方甘草片 3片 每日3次',
            advice: '注意休息，避免受凉，多喝温水，3天后复诊'
        },
        {
            hospital: '温江社区医疗中心',
            patientName: '周杰伦',
            diagnosis: '慢性胃炎',
            symptoms: '上腹部隐痛、餐后饱胀、反酸、嗳气',
            treatment: '药物治疗，饮食调理',
            medications: '奥美拉唑肠溶胶囊 20mg 每日1次；铝碳酸镁片 0.5g 每日3次；莫沙必利片 5mg 每日3次',
            advice: '规律饮食，避免辛辣刺激食物，戒烟限酒，按时服药，1个月后复查胃镜'
        },
        {
            hospital: '温江社区医疗中心',
            patientName: '周杰伦',
            diagnosis: '高血压（1级）',
            symptoms: '头晕、头痛、血压150/95mmHg',
            treatment: '降压药物治疗，生活方式干预',
            medications: '苯磺酸氨氯地平片 5mg 每日1次；缬沙坦胶囊 80mg 每日1次',
            advice: '低盐低脂饮食，适量运动，监测血压，每周复诊调整用药'
        },
        {
            hospital: '协和医院',
            patientName: '周杰伦',
            diagnosis: '2型糖尿病',
            symptoms: '多饮、多尿、体重下降、空腹血糖8.5mmol/L',
            treatment: '降糖药物治疗，饮食控制，运动疗法',
            medications: '二甲双胍缓释片 0.5g 每日2次；格列美脲片 2mg 每日1次；阿卡波糖片 50mg 每日3次',
            advice: '控制饮食总热量，少食多餐，规律运动，监测血糖，定期复查糖化血红蛋白'
        },
        {
            hospital: '协和医院',
            patientName: '周杰伦',
            diagnosis: '腰椎间盘突出症',
            symptoms: '腰痛伴左下肢放射痛、麻木，活动受限',
            treatment: '保守治疗，理疗，药物镇痛',
            medications: '塞来昔布胶囊 0.2g 每日1次；甲钴胺片 0.5mg 每日3次；盐酸乙哌立松片 50mg 每日3次',
            advice: '卧床休息，避免久坐久站，避免弯腰负重，配合理疗，2周后复诊评估'
        }
    ];

    console.log('🚀 开始为周杰伦创建5个病历...\n');

    for (let i = 0; i < prescriptions.length; i++) {
        const prescription = prescriptions[i];
        
        try {
            console.log(`📝 正在创建第 ${i + 1} 个病历 (${prescription.hospital})...`);
            
            const response = await fetch(`${baseURL}/api/v2/createPrescription`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    args: [{
                        patient_name: prescription.patientName,
                        diagnosis: prescription.diagnosis,
                        symptoms: prescription.symptoms,
                        treatment: prescription.treatment,
                        medications: prescription.medications,
                        advice: prescription.advice
                    }]
                })
            });

            const result = await response.json();
            
            if (response.ok && result.code === 200) {
                console.log(`✅ 第 ${i + 1} 个病历创建成功！`);
                console.log(`   医院: ${prescription.hospital}`);
                console.log(`   诊断: ${prescription.diagnosis}`);
                console.log(`   病历ID: ${result.data?.id || '未知'}\n`);
            } else {
                console.error(`❌ 第 ${i + 1} 个病历创建失败:`, result.message || '未知错误');
            }
            
            // 延迟500ms，避免请求过快
            await new Promise(resolve => setTimeout(resolve, 500));
            
        } catch (error) {
            console.error(`❌ 第 ${i + 1} 个病历创建出错:`, error.message);
        }
    }

    console.log('\n✨ 病历创建任务完成！');
    console.log('📊 统计：');
    console.log('   - 温江社区医疗中心: 3个病历');
    console.log('   - 协和医院: 2个病历');
    console.log('   - 总计: 5个病历');
}

// 执行函数
createPrescriptionsForZhouJielun();
